package handlers

import (
	"bytes"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/gabehf/koito/internal/catalog"
	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func ImageHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		size := chi.URLParam(r, "size")
		filename := chi.URLParam(r, "filename")

		imageSize, err := catalog.ParseImageSize(size)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		imgid, err := uuid.Parse(filename)
		if err != nil {
			serveDefaultImage(w, r, imageSize)
			return
		}

		desiredImgPath := filepath.Join(cfg.ConfigDir(), catalog.ImageCacheDir, size, filepath.Clean(filename))

		if _, err := os.Stat(desiredImgPath); os.IsNotExist(err) {
			l.Debug().Msg("Image not found in desired size")
			// file doesn't exist in desired size

			fullSizePath := filepath.Join(cfg.ConfigDir(), catalog.ImageCacheDir, string(catalog.ImageSizeFull), filepath.Clean(filename))
			largeSizePath := filepath.Join(cfg.ConfigDir(), catalog.ImageCacheDir, string(catalog.ImageSizeLarge), filepath.Clean(filename))

			// check if file exists at either full or large size
			// note: have to check both in case a user switched caching full size on and off
			// which would result in cache misses from source changing
			var sourcePath string
			if _, err = os.Stat(fullSizePath); os.IsNotExist(err) {
				if _, err = os.Stat(largeSizePath); os.IsNotExist(err) {
					l.Warn().Msgf("Could not find requested image %s. If this image is tied to an album or artist, it should be replaced", imgid.String())
					serveDefaultImage(w, r, imageSize)
					return
				} else if err != nil {
					// non-not found error for full file
					l.Err(err).Msg("Failed to access source image file")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				sourcePath = largeSizePath
			} else if err != nil {
				// non-not found error for full file
				l.Err(err).Msg("Failed to access source image file")
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				sourcePath = fullSizePath
			}

			// source size file was found

			// create and cache image at desired size

			imageBuf, err := os.ReadFile(sourcePath)
			if err != nil {
				l.Err(err).Msg("Failed to read source image file")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = catalog.CompressAndSaveImage(r.Context(), imgid.String(), imageSize, bytes.NewReader(imageBuf))
			if err != nil {
				l.Err(err).Msg("Failed to save compressed image to cache")
			}
		} else if err != nil {
			// non-not found error for desired file
			l.Err(err).Msg("Failed to access desired image file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Serve image
		http.ServeFile(w, r, desiredImgPath)
	}
}

func serveDefaultImage(w http.ResponseWriter, r *http.Request, size catalog.ImageSize) {
	var lock sync.Mutex
	l := logger.FromContext(r.Context())
	defaultImagePath := filepath.Join(cfg.ConfigDir(), catalog.ImageCacheDir, string(size), "default_img")
	if _, err := os.Stat(defaultImagePath); os.IsNotExist(err) {
		l.Debug().Msg("Default image does not exist in cache at desired size")
		defaultImagePath := filepath.Join(catalog.SourceImageDir(), "default_img")
		if _, err = os.Stat(defaultImagePath); os.IsNotExist(err) {
			l.Debug().Msg("Default image does not exist in cache, attempting to move...")
			err = os.MkdirAll(filepath.Dir(defaultImagePath), 0755)
			if err != nil {
				l.Err(err).Msg("Error when attempting to create image_cache/full dir")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			lock.Lock()
			utils.CopyFile(path.Join("assets", "default_img"), defaultImagePath)
			lock.Unlock()
		} else if err != nil {
			// non-not found error
			l.Error().Err(err).Msg("Error when attempting to read default image in cache")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// default_img does (or now does) exist in cache at full size
		file, err := os.Open(path.Join(catalog.SourceImageDir(), "default_img"))
		if err != nil {
			l.Err(err).Msg("Error when reading default image from source dir")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = catalog.CompressAndSaveImage(r.Context(), "default_img", size, file)
		if err != nil {
			l.Err(err).Msg("Error when caching default img at desired size")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		// non-not found error
		l.Error().Err(err).Msg("Error when attempting to read default image in cache")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// serve default_img at desired size
	http.ServeFile(w, r, path.Join(cfg.ConfigDir(), catalog.ImageCacheDir, string(size), "default_img"))
}

// func SearchMissingAlbumImagesHandler(store db.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		l := logger.FromContext(ctx)
// 		l.Info().Msg("Beginning search for albums with missing images")
// 		go func() {
// 			defer func() {
// 				if r := recover(); r != nil {
// 					l.Error().Interface("recover", r).Msg("Panic when searching for missing album images")
// 				}
// 			}()
// 			ctx := logger.NewContext(l)
// 			from := int32(0)
// 			count := 0
// 			for {
// 				albums, err := store.AlbumsWithoutImages(ctx, from)
// 				if errors.Is(err, pgx.ErrNoRows) {
// 					break
// 				} else if err != nil {
// 					l.Err(err).Msg("Failed to search for missing images")
// 					return
// 				}
// 				l.Debug().Msgf("Queried %d albums on page %d", len(albums), from)
// 				if len(albums) < 1 {
// 					break
// 				}
// 				for _, a := range albums {
// 					l.Debug().Msgf("Searching images for album %s", a.Title)
// 					img, err := imagesrc.GetAlbumImages(ctx, imagesrc.AlbumImageOpts{
// 						Artists:      utils.FlattenSimpleArtistNames(a.Artists),
// 						Album:        a.Title,
// 						ReleaseMbzID: a.MbzID,
// 					})
// 					if err == nil && img != "" {
// 						l.Debug().Msg("Image found! Downloading...")
// 						imgid, err := catalog.DownloadAndCacheImage(ctx, img)
// 						if err != nil {
// 							l.Err(err).Msgf("Failed to download image for %s", a.Title)
// 							continue
// 						}
// 						err = store.UpdateAlbum(ctx, db.UpdateAlbumOpts{
// 							ID:    a.ID,
// 							Image: imgid,
// 						})
// 						if err != nil {
// 							l.Err(err).Msgf("Failed to update image for %s", a.Title)
// 							continue
// 						}
// 						l.Info().Msgf("Found new album image for %s", a.Title)
// 						count++
// 					}
// 					if err != nil {
// 						l.Err(err).Msgf("Failed to get album images for %s", a.Title)
// 					}
// 				}
// 				from = albums[len(albums)-1].ID
// 			}
// 			l.Info().Msgf("Completed search, finding %d new images", count)
// 		}()
// 		w.WriteHeader(http.StatusOK)
// 	}
// }
