package repository

import (
	"golang.org/x/net/context"
	"soi-go/pkg/cli/config"
	"soi-go/pkg/cli/constant"
	"soi-go/pkg/model"
)

const (
	RepoTypeFile   = "file"
	RepoTypeSQLite = "sqlite"
)

type Repository interface {
	Init(ctx context.Context) error
	ListBucket(ctx context.Context) ([]string, error)
	LoadAll(ctx context.Context, bucket string) ([]*model.SoiData, error)
	Load(ctx context.Context, bucket string, hash string) (*model.SoiData, bool, error)
	Store(ctx context.Context, bucket string, soi *model.SoiData) error
	Exists(ctx context.Context, bucket string, hash string) (bool, error)
	Remove(ctx context.Context, bucket string, hash string) error
}

func NewRepository(ctx context.Context, conf *config.Config) (Repository, bool, error) {
	soisDir, err := constant.SoisDir()
	if err != nil {
		return nil, false, err
	}
	switch conf.DefaultRepository {
	case RepoTypeFile:
		repo, err := NewFilesRepository(soisDir)
		if err != nil {
			return nil, true, err
		}
		return repo, true, nil
	case RepoTypeSQLite:
		repo, err := NewSQLiteRepository(ctx, soisDir, conf.DefaultBucket)
		if err != nil {
			return nil, true, err
		}
		return repo, true, nil
	}
	return nil, false, nil
}
