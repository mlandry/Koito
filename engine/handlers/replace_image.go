package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gabehf/koito/internal/catalog"
	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
	"github.com/google/uuid"
)

type ReplaceImageResponse struct {
	Success bool   `json:"success"`
	Image   string `json:"image"`
	Message string `json:"message,omitempty"`
}

func ReplaceImageHandler(store db.DB, ip *catalog.ImageProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		artistIdStr := r.FormValue("artist_id")
		artistId, _ := strconv.Atoi(artistIdStr)
		albumIdStr := r.FormValue("album_id")
		albumId, _ := strconv.Atoi(albumIdStr)

		if artistId != 0 && albumId != 0 {
			utils.WriteError(w, "Only one of artist_id and album_id can be set", http.StatusBadRequest)
			return
		} else if artistId == 0 && albumId == 0 {
			utils.WriteError(w, "One of artist_id and album_id must be set", http.StatusBadRequest)
			return
		}

		var oldImage *uuid.UUID
		if artistId != 0 {
			a, err := store.GetArtist(ctx, db.GetArtistOpts{
				ID: int32(artistId),
			})
			if err != nil {
				utils.WriteError(w, "Artist with specified id could not be found", http.StatusBadRequest)
				return
			}
			oldImage = a.Image
		} else if albumId != 0 {
			a, err := store.GetAlbum(ctx, db.GetAlbumOpts{
				ID: int32(albumId),
			})
			if err != nil {
				utils.WriteError(w, "Album with specified id could not be found", http.StatusBadRequest)
				return
			}
			oldImage = a.Image
		}

		l.Debug().Msgf("Getting image from request...")

		var id uuid.UUID
		var err error

		fileUrl := r.FormValue("image_url")
		if fileUrl != "" {
			l.Debug().Msg("Image identified as remote file")
			err = catalog.ValidateImageURL(fileUrl)
			if err != nil {
				utils.WriteError(w, "url is invalid or not an image file", http.StatusBadRequest)
				return
			}
			id = uuid.New()
			var dlSize catalog.ImageSize
			if cfg.FullImageCacheEnabled() {
				dlSize = catalog.ImageSizeFull
			} else {
				dlSize = catalog.ImageSizeLarge
			}
			l.Debug().Msg("Downloading album image from source...")
			err = ip.EnqueueDownloadAndCache(ctx, id, fileUrl, dlSize)
			if err != nil {
				l.Err(err).Msg("Failed to cache image")
			}
		} else {
			file, _, err := r.FormFile("image")
			if err != nil {
				utils.WriteError(w, "Invalid file", http.StatusBadRequest)
				return
			}
			defer file.Close()

			buf := make([]byte, 512)
			if _, err := file.Read(buf); err != nil {
				utils.WriteError(w, "Could not read file", http.StatusInternalServerError)
				return
			}

			contentType := http.DetectContentType(buf)
			if !strings.HasPrefix(contentType, "image/") {
				utils.WriteError(w, "Only image uploads are allowed", http.StatusBadRequest)
				return
			}

			if _, err := file.Seek(0, io.SeekStart); err != nil {
				utils.WriteError(w, "Could not seek file", http.StatusInternalServerError)
				return
			}

			l.Debug().Msgf("Saving image to cache...")

			id = uuid.New()

			var dlSize catalog.ImageSize
			if cfg.FullImageCacheEnabled() {
				dlSize = catalog.ImageSizeFull
			} else {
				dlSize = catalog.ImageSizeLarge
			}

			err = ip.EnqueueCompressAndSave(ctx, id.String(), dlSize, file)
			if err != nil {
				utils.WriteError(w, "Could not save file", http.StatusInternalServerError)
				return
			}
		}

		l.Debug().Msgf("Updating database...")

		var imgsrc string
		if fileUrl != "" {
			imgsrc = fileUrl
		} else {
			imgsrc = catalog.ImageSourceUserUpload
		}

		if artistId != 0 {
			err := store.UpdateArtist(ctx, db.UpdateArtistOpts{
				ID:       int32(artistId),
				Image:    id,
				ImageSrc: imgsrc,
			})
			if err != nil {
				l.Err(err).Msg("Artist image could not be updated")
				utils.WriteError(w, "Artist image could not be updated", http.StatusInternalServerError)
				return
			}
		} else if albumId != 0 {
			err := store.UpdateAlbum(ctx, db.UpdateAlbumOpts{
				ID:       int32(albumId),
				Image:    id,
				ImageSrc: imgsrc,
			})
			if err != nil {
				l.Err(err).Msg("Album image could not be updated")
				utils.WriteError(w, "Album image could not be updated", http.StatusInternalServerError)
				return
			}
		}

		if oldImage != nil {
			l.Debug().Msg("Cleaning up old image file...")
			err = catalog.DeleteImage(*oldImage)
			if err != nil {
				l.Err(err).Msg("Failed to delete old image file")
				utils.WriteError(w, "Could not delete old image file", http.StatusInternalServerError)
				return
			}
		}

		utils.WriteJSON(w, http.StatusOK, ReplaceImageResponse{
			Success: true,
			Image:   id.String(),
		})
	}
}
