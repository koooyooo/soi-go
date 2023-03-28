package repository

import (
	"context"
	"github.com/koooyooo/soi-go/pkg/cli/loader"
	"github.com/koooyooo/soi-go/pkg/model"
	"path/filepath"
)

type jsonRepository struct {
	basePath string
}

func NewJsonsRepository(path string) (Repository, error) {
	return &jsonRepository{
		basePath: path,
	}, nil
}

func (r *jsonRepository) Init(ctx context.Context) error {
	return nil
}

func (r *jsonRepository) LoadAll(ctx context.Context, bucket string) ([]*model.SoiData, error) {
	return loader.LoadSois(filepath.Join(r.basePath, bucket))
}

func (r *jsonRepository) Load(ctx context.Context, bucket string, hash string) (*model.SoiData, bool, error) {
	sois, err := r.LoadAll(ctx, bucket)
	if err != nil {
		return nil, false, err
	}
	for _, soi := range sois {
		if soi.Hash == hash {
			return soi, true, nil
		}
	}
	return nil, false, nil
}

func (r *jsonRepository) Store(ctx context.Context, bucket string, soi *model.SoiData) error {
	return loader.StoreSoiData(filepath.Join(r.basePath, bucket, soi.Path+".json"), soi)
}

// TODO to path
func (r *jsonRepository) Exists(ctx context.Context, bucket string, hash string) (bool, error) {
	return loader.Exists(filepath.Join(r.basePath, bucket, hash+".json")), nil
}

// TODO to path
func (r *jsonRepository) Remove(ctx context.Context, bucket string, hash string) error {
	filepath.Join(r.basePath, bucket, hash+".json")
	return nil
}
