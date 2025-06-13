package engine

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gabehf/koito/engine/middleware"
	"github.com/gabehf/koito/internal/catalog"
	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/db/psql"
	"github.com/gabehf/koito/internal/images"
	"github.com/gabehf/koito/internal/importer"
	"github.com/gabehf/koito/internal/logger"
	mbz "github.com/gabehf/koito/internal/mbz"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/utils"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

const Version = "dev"

func Run(
	getenv func(string) string,
	w io.Writer,
) error {
	err := cfg.Load(getenv)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}
	l := logger.Get()
	if cfg.StructuredLogging() {
		*l = l.Output(w)
	} else {
		*l = l.Output(zerolog.ConsoleWriter{
			Out:        w,
			TimeFormat: time.RFC3339,
			// FormatLevel: func(i interface{}) string {
			// 	return strings.ToUpper(fmt.Sprintf("[%s]", i))
			// },
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("\u001b[30;1m>\u001b[0m %s |", i)
			},
		})
	}

	ctx := logger.NewContext(l)

	l.Info().Msgf("Koito %s", Version)

	_, err = os.Stat(cfg.ConfigDir())
	if err != nil {
		l.Info().Msgf("Creating config dir: %s", cfg.ConfigDir())
		err = os.MkdirAll(cfg.ConfigDir(), 0744)
		if err != nil {
			l.Error().Err(err).Msg("Failed to create config directory")
			return err
		}
	}
	l.Info().Msgf("Using config dir: %s", cfg.ConfigDir())
	_, err = os.Stat(path.Join(cfg.ConfigDir(), "import"))
	if err != nil {
		l.Debug().Msgf("Creating import dir: %s", path.Join(cfg.ConfigDir(), "import"))
		err = os.Mkdir(path.Join(cfg.ConfigDir(), "import"), 0744)
		if err != nil {
			l.Error().Err(err).Msg("Failed to create import directory")
			return err
		}
	}

	var store *psql.Psql
	store, err = psql.New()
	for err != nil {
		l.Error().Err(err).Msg("Failed to connect to database; retrying in 5 seconds")
		time.Sleep(5 * time.Second)
		store, err = psql.New()
	}
	defer store.Close(ctx)

	var mbzC mbz.MusicBrainzCaller
	if !cfg.MusicBrainzDisabled() {
		mbzC = mbz.NewMusicBrainzClient()
	} else {
		mbzC = &mbz.MbzErrorCaller{}
	}

	images.Initialize(images.ImageSourceOpts{
		UserAgent:    cfg.UserAgent(),
		EnableCAA:    !cfg.CoverArtArchiveDisabled(),
		EnableDeezer: !cfg.DeezerDisabled(),
	})

	userCount, _ := store.CountUsers(ctx)
	if userCount < 1 {
		l.Debug().Msg("Creating default user...")
		user, err := store.SaveUser(ctx, db.SaveUserOpts{
			Username: cfg.DefaultUsername(),
			Password: cfg.DefaultPassword(),
			Role:     models.UserRoleAdmin,
		})
		if err != nil {
			l.Fatal().AnErr("error", err).Msg("Failed to save default user in database")
		}
		apikey, err := utils.GenerateRandomString(48)
		if err != nil {
			l.Fatal().AnErr("error", err).Msg("Failed to generate default api key")
		}
		label := "Default"
		_, err = store.SaveApiKey(ctx, db.SaveApiKeyOpts{
			Key:    apikey,
			UserID: user.ID,
			Label:  label,
		})
		if err != nil {
			l.Fatal().AnErr("error", err).Msg("Failed to save default api key in database")
		}
		l.Info().Msgf("Default user has been created. Login: %s : %s", cfg.DefaultUsername(), cfg.DefaultPassword())
	}

	if cfg.AllowAllHosts() {
		l.Warn().Msg("Your configuration allows requests from all hosts. This is a potential security risk!")
	} else if len(cfg.AllowedHosts()) == 0 || cfg.AllowedHosts()[0] == "" {
		l.Warn().Msgf("You are currently not allowing any hosts! Did you forget to set the %s variable?", cfg.ALLOWED_HOSTS_ENV)
	} else {
		l.Debug().Msgf("Allowing hosts: %v", cfg.AllowedHosts())
	}

	var ready atomic.Bool

	mux := chi.NewRouter()
	// bind general middleware to mux
	mux.Use(middleware.WithRequestID)
	mux.Use(middleware.Logger(l))
	mux.Use(chimiddleware.Recoverer)
	mux.Use(chimiddleware.RealIP)
	// call router binds on mux
	bindRoutes(mux, &ready, store, mbzC)

	httpServer := &http.Server{
		Addr:    cfg.ListenAddr(),
		Handler: mux,
	}

	go func() {
		ready.Store(true) // signal readiness
		l.Info().Msg("listening on " + cfg.ListenAddr())
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal().AnErr("error", err).Msg("Error when running ListenAndServe")
		}
	}()

	// Import
	if !cfg.SkipImport() {
		go func() {
			RunImporter(l, store, mbzC)
		}()
	}

	l.Info().Msg("Pruning orphaned images...")
	go catalog.PruneOrphanedImages(logger.NewContext(l), store)
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	l.Info().Msg("Received server shutdown notice")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	l.Info().Msg("waiting for all processes to finish...")
	mbzC.Shutdown()
	if err := httpServer.Shutdown(ctx); err != nil {
		return err
	}
	l.Info().Msg("shutdown successful")
	return nil
}

func RunImporter(l *zerolog.Logger, store db.DB, mbzc mbz.MusicBrainzCaller) {
	l.Debug().Msg("Checking for import files...")
	files, err := os.ReadDir(path.Join(cfg.ConfigDir(), "import"))
	if err != nil {
		l.Err(err).Msg("Failed to read files from import dir")
	}
	if len(files) > 0 {
		l.Info().Msg("Files found in import directory. Attempting to import...")
	} else {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			l.Error().Interface("recover", r).Msg("Panic when importing files")
		}
	}()
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), "Streaming_History_Audio") {
			l.Info().Msgf("Import file %s detecting as being Spotify export", file.Name())
			err := importer.ImportSpotifyFile(logger.NewContext(l), store, file.Name())
			if err != nil {
				l.Err(err).Msgf("Failed to import file: %s", file.Name())
			}
		} else if strings.Contains(file.Name(), "maloja") {
			l.Info().Msgf("Import file %s detecting as being Maloja export", file.Name())
			err := importer.ImportMalojaFile(logger.NewContext(l), store, file.Name())
			if err != nil {
				l.Err(err).Msgf("Failed to import file: %s", file.Name())
			}
		} else if strings.Contains(file.Name(), "recenttracks") {
			l.Info().Msgf("Import file %s detecting as being ghan.nl LastFM export", file.Name())
			err := importer.ImportLastFMFile(logger.NewContext(l), store, mbzc, file.Name())
			if err != nil {
				l.Err(err).Msgf("Failed to import file: %s", file.Name())
			}
		} else {
			l.Warn().Msgf("File %s not recognized as a valid import file; make sure it is valid and named correctly", file.Name())
		}
	}
}
