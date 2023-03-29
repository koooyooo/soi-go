package service

import (
	"context"
	"github.com/koooyooo/soi-go/pkg/cli/repository"
	"github.com/koooyooo/soi-go/pkg/common/hash"
	"github.com/koooyooo/soi-go/pkg/model"
	"strings"
)

type Service interface {
	Init(ctx context.Context) error
	ListBucket(ctx context.Context) ([]string, error)
	ChangeBucket(ctx context.Context, bucket string) error
	LoadAll(ctx context.Context) ([]*model.SoiData, error)
	Load(ctx context.Context, hash string) (*model.SoiData, bool, error)
	Store(ctx context.Context, soi *model.SoiData) error
	Exists(ctx context.Context, hash string) (bool, error)
	Remove(ctx context.Context, hash string) error
	ListDir(ctx context.Context) ([]string, error)
}

func NewService(ctx context.Context, bucket string, r repository.Repository) Service {
	return &serviceImpl{
		bucket: bucket,
		r:      r,
	}
}

type serviceImpl struct {
	bucket string
	r      repository.Repository
}

func (s serviceImpl) Init(ctx context.Context) error {
	return s.r.Init(ctx)
}

func (s serviceImpl) ListBucket(ctx context.Context) ([]string, error) {
	return s.r.ListBucket(ctx)
}

func (s serviceImpl) ChangeBucket(ctx context.Context, bucket string) error {
	s.bucket = bucket
	return nil
}

func (s serviceImpl) LoadAll(ctx context.Context) ([]*model.SoiData, error) {
	return s.r.LoadAll(ctx, s.bucket)
}

func (s serviceImpl) Load(ctx context.Context, hash string) (*model.SoiData, bool, error) {
	return s.r.Load(ctx, s.bucket, hash)
}

func (s serviceImpl) Store(ctx context.Context, soi *model.SoiData) error {
	return s.r.Store(ctx, s.bucket, soi)
}

func (s serviceImpl) Exists(ctx context.Context, hash string) (bool, error) {
	return s.r.Exists(ctx, s.bucket, hash)
}

func (s serviceImpl) Remove(ctx context.Context, hash string) error {
	return s.r.Remove(ctx, s.bucket, hash)
}

func (s serviceImpl) ListDir(ctx context.Context) ([]string, error) {
	sois, err := s.LoadAll(ctx)
	if err != nil {
		return nil, err
	}
	var m map[string]struct{}
	for _, soi := range sois {
		m[soi.Path] = struct{}{}
	}
	var dirs []string
	for k, _ := range m {
		dirs = append(dirs, k)
	}
	return dirs, nil
}

func toHash(path string) (string, error) {
	return hash.Sha1(path)
}

func findHashes(r repository.Repository, bucket, partialPath string) ([]string, error) {
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