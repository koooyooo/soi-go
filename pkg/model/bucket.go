package model

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/koooyooo/soi-go/pkg/constant"
)

type BucketRef struct {
	Bucket *Bucket
}

func NewBucket(name string) (*Bucket, error) {
	baseDir, err := constant.SoisDir()
	if err != nil {
		return nil, err
	}
	return &Bucket{
		BaseDir: baseDir,
		Name:    name,
	}, nil
}

// Bucket はローカルバケットを管理する構造体
type Bucket struct {
	BaseDir string
	Name    string
}

// IsLocalOnly はローカル限定のバケット
func (l Bucket) IsLocalOnly() bool {
	return strings.HasPrefix(l.Name, "_")
}

// Path はローカルバケット毎のルートパスを取得する
func (l Bucket) Path() (string, error) {
	return filepath.Join(l.BaseDir, l.Name), nil
}

func ListBuckets() ([]*Bucket, error) {
	var buckets []*Bucket
	bucketNames, err := ListBucketNames()
	if err != nil {
		return nil, err
	}
	for _, n := range bucketNames {
		b, err := NewBucket(n)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, b)
	}
	return buckets, nil
}

// listBuckets はバケット一覧を取得する
func ListBucketNames() ([]string, error) {
	soisDir, err := constant.SoisDir()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(soisDir)
	if err != nil {
		return nil, err
	}

	var buckets []string
	for _, f := range files {
		if f.IsDir() {
			buckets = append(buckets, f.Name())
		}
	}
	return buckets, nil
}
