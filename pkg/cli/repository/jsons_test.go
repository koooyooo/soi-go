package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONsLoadAll(t *testing.T) {
	repo, err := NewJsonsRepository("../../../testfiles")
	assert.NoError(t, err)

	ctx := context.Background()
	sois, err := repo.LoadAll(ctx, "bucket1")
	assert.NoError(t, err)
	assert.Equal(t, "google", sois[0].Name)
}

func TestJSONsLoad(t *testing.T) {
	repo, err := NewJsonsRepository("../../../testfiles")
	assert.NoError(t, err)

	ctx := context.Background()
	soi, ok, err := repo.Load(ctx, "bucket1", "aa4f20c99b3c3d188ce5b6255a299a32d7fa9c78")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "google", soi.Name)
}

func TestJSONsStore(t *testing.T) {
	// TODO
}
