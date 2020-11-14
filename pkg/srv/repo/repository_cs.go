package repo

import (
	"context"
	"encoding/json"
	"path/filepath"

	"cloud.google.com/go/storage"

	"github.com/koooyooo/soi-go/pkg/soi"
)

type CloudStoreRepository struct {
	BucketName string
}

func (r CloudStoreRepository) StoreAll(ctx context.Context, sv *soi.SoiVirtual) error {
	bkt, err := getBucket(ctx, r.BucketName)
	if err != nil {
		return err
	}
	objName, err := getObjectName(ctx)
	if err != nil {
		return err
	}
	obj := bkt.Object(objName)
	w := obj.NewWriter(ctx)
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

func getBucket(ctx context.Context, bucketName string) (*storage.BucketHandle, error) {
	cli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return cli.Bucket(bucketName), nil
}

func getObjectName(ctx context.Context) (string, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return "", err
	}
	soiBucketID, err := getSoiBucketID(ctx)
	if err != nil {
		return "", err
	}
	return filepath.Join(userID, soiBucketID+".json"), nil
}
