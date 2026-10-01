package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/koooyooo/soi-go/pkg/model"
)

func TestJSONsLoadAll(t *testing.T) {
	repo, err := NewFilesRepository("../../testfiles")
	assert.NoError(t, err)

	ctx := context.Background()
	sois, err := repo.LoadAll(ctx, "bucket1")
	assert.NoError(t, err)
	require.NotEmpty(t, sois)

	names := map[string]bool{}
	for _, s := range sois {
		names[s.Name] = true
		assert.Equal(t, "portal", s.Path)
	}
	assert.True(t, names["google"])
}

func TestJSONsLoad(t *testing.T) {
	repo, err := NewFilesRepository("../../testfiles")
	assert.NoError(t, err)

	ctx := context.Background()
	soi, ok, err := repo.Load(ctx, "bucket1", "aa4f20c99b3c3d188ce5b6255a299a32d7fa9c78")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "google", soi.Name)
	assert.Equal(t, "portal", soi.Path)
}

func TestJSONsStoreAndRemove(t *testing.T) {
	base := t.TempDir()
	repo, err := NewFilesRepository(base)
	require.NoError(t, err)

	ctx := context.Background()
	soi := &model.SoiData{
		Name:      "example",
		Path:      "docs",
		URI:       "https://example.com",
		Hash:      "testhash123",
		CreatedAt: time.Now(),
	}
	require.NoError(t, repo.Store(ctx, "bucket1", soi))

	abs := filepath.Join(base, soi.FilePath("bucket1"))
	_, err = os.Stat(abs)
	require.NoError(t, err)

	loaded, ok, err := repo.Load(ctx, "bucket1", "testhash123")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "example", loaded.Name)
	assert.Equal(t, "docs", loaded.Path)

	require.NoError(t, repo.Remove(ctx, "bucket1", "testhash123"))
	_, err = os.Stat(abs)
	assert.True(t, os.IsNotExist(err))
}
