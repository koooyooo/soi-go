package service

import (
	"context"
	"sort"
	"strings"

	"github.com/koooyooo/soi-go/pkg/common/hash"
	"github.com/koooyooo/soi-go/pkg/model"
	"github.com/koooyooo/soi-go/pkg/repository"
)

type Service interface {
	Init(ctx context.Context) error
	ListBucket(ctx context.Context) ([]string, error)
	CurrentBucket(ctx context.Context) (string, error)
	ChangeBucket(ctx context.Context, bucket string) error
	LoadAll(ctx context.Context) ([]*model.SoiData, error)
	Load(ctx context.Context, hash string) (*model.SoiData, bool, error)
	Store(ctx context.Context, soi *model.SoiData) error
	Exists(ctx context.Context, hash string) (bool, error)
	Remove(ctx context.Context, hash string) error
	ListPath(ctx context.Context, partialPath string, withName bool) ([]string, error)
	Size(ctx context.Context) (int, error)
}

type serviceImpl struct {
	bucket string
	r      repository.Repository
}

func NewService(_ context.Context, bucket string, r repository.Repository) Service {
	return &serviceImpl{
		bucket: bucket,
		r:      r,
	}
}

func (s *serviceImpl) Init(ctx context.Context) error {
	return s.r.Init(ctx)
}

func (s *serviceImpl) ListBucket(ctx context.Context) ([]string, error) {
	return s.r.ListBucket(ctx)
}

func (s *serviceImpl) CurrentBucket(ctx context.Context) (string, error) {
	return s.bucket, nil
}

func (s *serviceImpl) ChangeBucket(_ context.Context, bucket string) error {
	s.bucket = bucket
	return nil
}

func (s *serviceImpl) LoadAll(ctx context.Context) ([]*model.SoiData, error) {
	return s.r.LoadAll(ctx, s.bucket)
}

func (s *serviceImpl) Load(ctx context.Context, hash string) (*model.SoiData, bool, error) {
	return s.r.Load(ctx, s.bucket, hash)
}

func (s *serviceImpl) Store(ctx context.Context, soi *model.SoiData) error {
	return s.r.Store(ctx, s.bucket, soi)
}

func (s *serviceImpl) Exists(ctx context.Context, hash string) (bool, error) {
	return s.r.Exists(ctx, s.bucket, hash)
}

func (s *serviceImpl) Remove(ctx context.Context, hash string) error {
	return s.r.Remove(ctx, s.bucket, hash)
}

func (s *serviceImpl) ListPath(ctx context.Context, partialPath string, withName bool) ([]string, error) {
	sois, err := s.LoadAll(ctx)
	if err != nil {
		return nil, err
	}
	var m = make(map[string]struct{})
	for _, soi := range sois {
		path := soi.Path
		if withName {
			path = path + "/" + soi.Name
		}
		if !strings.HasPrefix(path, partialPath) {
			continue
		}
		elm := strings.Split(path, "/")
		for i := 0; i < len(elm); i++ {
			path := strings.Join(elm[:i+1], "/")
			if path == "" {
				continue
			}
			if !withName || i != len(elm)-1 {
				path = path + "/"
			}
			m[path] = struct{}{}
		}
	}
	var dirs []string
	for k, _ := range m {
		if !strings.HasPrefix(k, partialPath) {
			continue
		}
		dirs = append(dirs, k)
	}
	sort.Strings(dirs)
	return dirs, nil
}

func (s *serviceImpl) Size(ctx context.Context) (int, error) {
	sois, err := s.LoadAll(ctx)
	if err != nil {
		return -1, err
	}
	return len(sois), nil
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
