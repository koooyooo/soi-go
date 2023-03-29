package repository

import (
	"context"
	"fmt"
	"github.com/koooyooo/soi-go/pkg/model"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestSQLite(t *testing.T) {
	os.RemoveAll("test.db")
	ctx := context.Background()
	repo, err := NewSQLiteRepository(ctx, "", "test")
	assert.NoError(t, err)
	err = repo.Init(context.Background())
	assert.NoError(t, err)
	fmt.Println(repo)
}

func TestStore(t *testing.T) {
	soi := &model.SoiData{
		Hash:  "test-hash",
		Name:  "test-name",
		Title: "test-title",
		Path:  "test-basePath",
		URI:   "test-uri",
		Tags:  []string{"tag1", "tag2", "tag3"},
		KVTags: []model.KVTag{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
		},
		UsageLogs: []model.UsageLog{
			{Type: model.UsageTypeOpen, UsedAt: time.Now()},
			{Type: model.UsageTypeOpen, UsedAt: time.Now()},
		},
		Rate:          0.5,
		OGTitle:       "test-og-title",
		OGURL:         "/og/url",
		OGType:        "og-type",
		OGDescription: "og-description",
		OGSiteName:    "OG-site name",
		OGImages: []model.OGImage{
			{
				URL:       "https://example.com/og-image1.png",
				SecureURL: "https://example.com/og-image1.png",
				Type:      "image/png",
				Width:     100,
				Height:    100,
				Alt:       "og-image1",
			},
			{
				URL:       "https://example.com/og-image2.png",
				SecureURL: "https://example.com/og-image2.png",
				Type:      "image/png",
				Width:     200,
				Height:    200,
				Alt:       "og-image2",
			},
			{
				URL:       "https://example.com/og-image3.png",
				SecureURL: "https://example.com/og-image3.png",
				Type:      "image/png",
				Width:     300,
				Height:    300,
				Alt:       "og-image3",
			},
		},
	}
	ctx := context.Background()

	os.RemoveAll("test.db")
	repo, err := NewSQLiteRepository(ctx, "", "test")
	assert.NoError(t, err)

	err = repo.Init(ctx)
	assert.NoError(t, err)

	err = repo.Store(ctx, "test", soi)
	assert.NoError(t, err)

	result, ok, err := repo.Load(context.Background(), "test", soi.Hash)
	assert.NoError(t, err)
	assert.True(t, ok)

	// check each attribute
	assert.Equal(t, soi.Hash, result.Hash)
	assert.Equal(t, soi.Name, result.Name)
	assert.Equal(t, soi.Title, result.Title)
	assert.Equal(t, soi.Path, result.Path)
	assert.Equal(t, soi.URI, result.URI)
	assert.Equal(t, soi.Tags, result.Tags)
	assert.Equal(t, soi.KVTags, result.KVTags)
	assert.Equal(t, soi.Rate, result.Rate)

	// og attribute
	assert.Equal(t, soi.OGTitle, result.OGTitle)
	assert.Equal(t, soi.OGURL, result.OGURL)
	assert.Equal(t, soi.OGType, result.OGType)
	assert.Equal(t, soi.OGDescription, result.OGDescription)
	assert.Equal(t, soi.OGSiteName, result.OGSiteName)

	// og image
	assert.Equal(t, len(soi.OGImages), len(result.OGImages))
	for i, img := range soi.OGImages {
		assert.Equal(t, img, result.OGImages[i])
	}

	for _, log := range result.UsageLogs {
		fmt.Println("type", log.Type, "at", log.UsedAt)
	}
}
