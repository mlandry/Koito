package catalog_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gabehf/koito/internal/catalog"
	"github.com/gabehf/koito/internal/cfg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageLifecycle(t *testing.T) {

	ip := catalog.NewImageProcessor(1)

	// serve yuu.jpg as test image
	imageBytes, err := os.ReadFile(filepath.Join("static", "yuu.jpg"))
	require.NoError(t, err)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		w.Write(imageBytes)
	}))
	defer server.Close()

	imgID := uuid.New()

	err = ip.EnqueueDownloadAndCache(context.Background(), imgID, server.URL, catalog.ImageSizeFull)
	require.NoError(t, err)
	err = ip.EnqueueDownloadAndCache(context.Background(), imgID, server.URL, catalog.ImageSizeMedium)
	require.NoError(t, err)

	ip.WaitForIdle(5 * time.Second)

	// ensure download is correct

	imagePath := filepath.Join(cfg.ConfigDir(), catalog.ImageCacheDir, "full", imgID.String())
	assert.NoError(t, waitForFile(imagePath, 1*time.Second))
	imagePath = filepath.Join(cfg.ConfigDir(), catalog.ImageCacheDir, "medium", imgID.String())
	assert.NoError(t, waitForFile(imagePath, 1*time.Second))

	assert.NoError(t, catalog.DeleteImage(imgID))

	// ensure delete works

	imagePath = filepath.Join(cfg.ConfigDir(), catalog.ImageCacheDir, "full", imgID.String())
	assert.Error(t, waitForFile(imagePath, 1*time.Second))
	imagePath = filepath.Join(cfg.ConfigDir(), catalog.ImageCacheDir, "medium", imgID.String())
	assert.Error(t, waitForFile(imagePath, 1*time.Second))

	// re-download for prune

	err = ip.EnqueueDownloadAndCache(context.Background(), imgID, server.URL, catalog.ImageSizeFull)
	require.NoError(t, err)
	err = ip.EnqueueDownloadAndCache(context.Background(), imgID, server.URL, catalog.ImageSizeMedium)
	require.NoError(t, err)

	ip.WaitForIdle(5 * time.Second)

	assert.NoError(t, catalog.PruneOrphanedImages(context.Background(), store))

	// ensure prune works

	imagePath = filepath.Join(cfg.ConfigDir(), catalog.ImageCacheDir, "full", imgID.String())
	assert.Error(t, waitForFile(imagePath, 1*time.Second))
	imagePath = filepath.Join(cfg.ConfigDir(), catalog.ImageCacheDir, "medium", imgID.String())
	assert.Error(t, waitForFile(imagePath, 1*time.Second))
}

func waitForFile(path string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		if _, err := os.Stat(path); err == nil {
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for %s", path)
		}
		time.Sleep(20 * time.Millisecond)
	}
}
