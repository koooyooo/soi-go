package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/koooyooo/soi-go/pkg/fileio"
	"github.com/koooyooo/soi-go/pkg/soi"
)

func newFileRepository() FileRepository {
	return FileRepository{
		BasePath: "./repo/",
	}
}

type FileRepository struct {
	BasePath string
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
	soiBucketID, err := getSoiBucketID(ctx)
	if err != nil {
		return err
	}
	p := path.Join(f.BasePath, userID, soiBucketID)
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
	soiBucketID, err := getSoiBucketID(ctx)
	if err != nil {
		return nil, err
	}
	path := path.Join(f.BasePath, userID, soiBucketID, "sois.json")
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
