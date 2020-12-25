package repo

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"cloud.google.com/go/storage"

	"github.com/koooyooo/soi-go/pkg/soi"
)

// GCSリポジトリを作成
func newGCSRepository(bucketName string) *GCSRepository {
	return &GCSRepository{BucketName: bucketName}
}

// GCSリポジトリ
type GCSRepository struct {
	// GCSバケット名
	BucketName string
}

// Store は単一の仮想Soiを保存する
func (gr GCSRepository) Store(ctx context.Context, sv *soi.SoiVirtual) error {
	return store(gr, ctx, sv)
}

// StoreAll は仮想Soiバケットを保存する
func (gr GCSRepository) StoreAll(ctx context.Context, sv *soi.SoiVirtualBucket) error {
	bkt, err := getBucket(ctx, gr.BucketName)
	if err != nil {
		return err
	}
	objName, err := getObjectName(ctx)
	if err != nil {
		return err
	}
	w := bkt.Object(objName).NewWriter(ctx)
	defer w.Close()
	b, err := json.Marshal(sv)

	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

// LoadAll は仮想Soiバケットを取得する
func (gr GCSRepository) LoadAll(ctx context.Context) (*soi.SoiVirtualBucket, error) {
	bkt, err := getBucket(ctx, gr.BucketName)
	if err != nil {
		return nil, err
	}
	objName, err := getObjectName(ctx)
	if err != nil {
		return nil, err
	}
	r, err := bkt.Object(objName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var svb soi.SoiVirtualBucket
	if err := json.Unmarshal(b, &svb); err != nil {
		return nil, err
	}
	return &svb, nil
}

// getBucket はContextよりバケット名を取得する
func getBucket(ctx context.Context, bucketName string) (*storage.BucketHandle, error) {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return cli.Bucket(bucketName), nil
}

// getObjectName はContextよりオブジェクト名を取得する
func getObjectName(ctx context.Context) (string, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return "", err
	}
	soiBucketID, err := getSoiBucketID(ctx)
	if err != nil {
		return "", err
	}
	return filepath.Join(userID, soiBucketID, "sois.json"), nil
}
