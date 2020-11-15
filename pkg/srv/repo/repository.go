package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/koooyooo/soi-go/pkg/srv/constant"

	"github.com/koooyooo/soi-go/pkg/soi"
)

type (
	Repository interface {
		Store(context.Context, *soi.SoiVirtual) error
		StoreAll(context.Context, *soi.SoiVirtualBucket) error
		LoadAll(context.Context) (*soi.SoiVirtualBucket, error)
	}
)

func NewRepository() Repository {
	repoType := 2
	switch repoType {
	case 1:
		return newFileRepository("./repo/")
	case 2:
		return newGCSRepository(constant.EnvKeySoiBucketName.Get())
	}
	return nil
}

func getUserID(ctx context.Context) (string, error) {
	return getValue(ctx, constant.CtxKeyUserID)
}

func getSoiBucketID(ctx context.Context) (string, error) {
	return getValue(ctx, constant.CtxKeySoiBucketID)
}

func getValue(ctx context.Context, key constant.CtxKey) (string, error) {
	if strings.Contains(key.String(), "..") {
		return "", fmt.Errorf("invalid key: %v", key)
	}
	val, ok := ctx.Value(key).(string)
	if !ok {
		return "", fmt.Errorf("no value found for key: %v", key)
	}
	return val, nil
}

func store(r Repository, ctx context.Context, s *soi.SoiVirtual) error {
	sb, err := r.LoadAll(ctx)
	if err != nil {
		return nil
	}
	sb.Sois = append(sb.Sois, s)
	return r.StoreAll(ctx, sb)
}
