package service_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/koooyooo/soi-go/pkg/model"
	"github.com/koooyooo/soi-go/pkg/repository"
	"github.com/koooyooo/soi-go/pkg/service"
)

func TestServiceStoreLoadRoundTrip(t *testing.T) {
	base := t.TempDir()
	repo, err := repository.NewFilesRepository(base)
	require.NoError(t, err)

	ctx := context.Background()
	svc := service.NewService(ctx, "default", repo)

	soi := &model.SoiData{
		Name:      "Round Trip",
		Path:      "notes",
		URI:       "https://example.com/round-trip",
		Hash:      "roundtriphash",
		Tags:      []string{"demo"},
		CreatedAt: time.Now(),
	}
	require.NoError(t, svc.Store(ctx, soi))

	loaded, ok, err := svc.Load(ctx, "roundtriphash")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "Round Trip", loaded.Name)
	require.Equal(t, "notes", loaded.Path)

	abs := filepath.Join(base, soi.FilePath("default"))
	b, err := os.ReadFile(abs)
	require.NoError(t, err)
	var disk model.SoiData
	require.NoError(t, json.Unmarshal(b, &disk))
	require.Equal(t, "https://example.com/round-trip", disk.URI)

	require.NoError(t, svc.Remove(ctx, "roundtriphash"))
	_, ok, err = svc.Load(ctx, "roundtriphash")
	require.NoError(t, err)
	require.False(t, ok)
}
