package importer

import (
	"context"
	"os"
	"path"

	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/logger"
)

func finishImport(ctx context.Context, filename string, numImported int) error {
	l := logger.FromContext(ctx)
	_, err := os.Stat(path.Join(cfg.ConfigDir(), "import_complete"))
	if err != nil {
		err = os.Mkdir(path.Join(cfg.ConfigDir(), "import_complete"), 0744)
		if err != nil {
			l.Err(err).Msg("Failed to create import_complete dir! Import files must be removed from the import directory manually, or else the importer will run on every app start")
		}
	}
	err = os.Rename(path.Join(cfg.ConfigDir(), "import", filename), path.Join(cfg.ConfigDir(), "import_complete", filename))
	if err != nil {
		l.Err(err).Msg("Failed to move file to import_complete dir! Import files must be removed from the import directory manually, or else the importer will run on every app start")
	}
	if numImported != 0 {
		l.Info().Msgf("Finished importing %s; imported %d items", filename, numImported)
	}
	return nil
}
