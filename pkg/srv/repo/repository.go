package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/koooyooo/soi-go/pkg/soi"

	"github.com/koooyooo/soi-go/pkg/srv"

	"github.com/koooyooo/soi-go/pkg/fileio"
)

type (
	Repository interface {
		Store(context.Context, *soi.SoiVirtual) error
		StoreAll(context.Context, *soi.SoiVirtualBucket) error
		LoadAll(context.Context) (*soi.SoiVirtualBucket, error)
	}

	FileRepository struct {
		BasePath string
	}
)

func NewRepository() Repository {
	return FileRepository{
		BasePath: "./repo/",
	}
}

func (f FileRepository) Store(ctx context.Context, s *soi.SoiVirtual) error {
	sb, err := f.LoadAll(ctx)
	if err != nil {
		return nil
	}
	sb.Sois = append(sb.Sois, s)
	return f.StoreAll(ctx, sb)
}

func (f FileRepository) StoreAll(ctx context.Context, sb *soi.SoiVirtualBucket) error {
	var buff bytes.Buffer
	buff.WriteString("{\"sois\":[\n")
	for i, sv := range sb.Sois {
		svb, err := json.Marshal(sv)
		if err != nil {
			return err
		}
		buff.WriteString("  " + string(svb))
		if i != len(sb.Sois)-1 {
			buff.WriteString(",")
		}
		buff.WriteString("\n")
	}
	buff.WriteString("]}\n")

	userID, err := getUserID(ctx)
	if err != nil {
		return err
	}
	p := path.Join(f.BasePath, userID)
	if err = os.MkdirAll(p, 0700); err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(p, "sois.json"), buff.Bytes(), 0600)
}

func (f FileRepository) LoadAll(ctx context.Context) (*soi.SoiVirtualBucket, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}
	path := path.Join(f.BasePath, userID, "sois.json")
	if !fileio.Exists(path) {
		return &soi.SoiVirtualBucket{}, nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var soiBucket soi.SoiVirtualBucket
	if err := json.Unmarshal(b, &soiBucket); err != nil {
		return nil, err
	}
	return &soiBucket, nil
}

func getUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(srv.CtxKeyUserID).(string)
	if !ok {
		return "", fmt.Errorf("no user found: %v", userID)
	}
	if strings.Contains(userID, "..") {
		return "", fmt.Errorf("invalid user id: %v", userID)
	}
	return userID, nil
}
