package repository

import (
	"context"
	"github.com/koooyooo/soi-go/pkg/cli/loader"
	"github.com/koooyooo/soi-go/pkg/model"
	"os"
	"path/filepath"
)

type filesRepository struct {
	basePath string
}

func NewFilesRepository(path string) (Repository, error) {
	return &filesRepository{
		basePath: path,
	}, nil
}

func (r *filesRepository) Init(ctx context.Context) error {
	return nil
}

func (r *filesRepository) LoadAll(ctx context.Context, bucket string) ([]*model.SoiData, error) {
	return loader.LoadSois(filepath.Join(r.basePath, bucket))
}

func (r *filesRepository) Load(ctx context.Context, bucket string, hash string) (*model.SoiData, bool, error) {
	sois, err := r.LoadAll(ctx, bucket)
	if err != nil {
		return nil, false, err
	}
	s, found := findSoi(sois, hash)
	return s, found, nil
}

func (r *filesRepository) Store(ctx context.Context, bucket string, soi *model.SoiData) error {
	return loader.StoreSoiData(soi.FilePath(bucket), soi)
}

func (r *filesRepository) Exists(ctx context.Context, bucket string, hash string) (bool, error) {
	_, found, err := r.Load(ctx, bucket, hash)
	return found, err
}

func (r *filesRepository) Remove(ctx context.Context, bucket string, hash string) error {
	s, found, err := r.Load(ctx, bucket, hash)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	return os.RemoveAll(s.FilePath(bucket))
}

func findSoi(sois []*model.SoiData, hash string) (*model.SoiData, bool) {
	for _, soi := range sois {
		if soi.Hash == hash {
			return soi, true
		}
	}
	return nil, false
}
