package repository

import (
	"github.com/koooyooo/soi-go/pkg/common/hash"
	"github.com/koooyooo/soi-go/pkg/model"
	"golang.org/x/net/context"
	"strings"
)

type Repository interface {
	Init(ctx context.Context) error
	// TODO add "limit" args and -1 means no limit
	LoadAll(ctx context.Context, bucket string) ([]*model.SoiData, error)
	Load(ctx context.Context, bucket string, hash string) (*model.SoiData, bool, error)
	Store(ctx context.Context, bucket string, soi *model.SoiData) error
	Exists(ctx context.Context, bucket string, hash string) (bool, error)
	//Remove(ctx context.Context, bucket string, hash string) error
}

func toHash(path string) (string, error) {
	return hash.Sha1(path)
}

func findHashes(r Repository, bucket, partialPath string) ([]string, error) {
	sois, err := r.LoadAll(context.Background(), bucket)
	if err != nil {
		return nil, err
	}
	var hashes []string
	for _, soi := range sois {
		if strings.Contains(soi.Path, partialPath) {
			hashes = append(hashes, soi.Hash)
		}
	}
	return hashes, nil
}
