package repository

import (
	"github.com/koooyooo/soi-go/pkg/common/hash"
	"github.com/koooyooo/soi-go/pkg/model"
	"golang.org/x/net/context"
)

type Repository interface {
	Init(ctx context.Context) error
	// TODO add "limit" args and -1 means no limit
	LoadAll(ctx context.Context, bucket string) ([]*model.SoiData, error)
	Load(ctx context.Context, bucket string, hash string) (*model.SoiData, bool, error)
	Store(ctx context.Context, bucket string, soi *model.SoiData) error
	// to path
	Exists(ctx context.Context, bucket string, hash string) (bool, error)
	// to path
	//Remove(ctx context.Context, bucket string, hash string) error
}

func toHash(path string) (string, error) {
	return hash.Sha1(path)
}
