package repo

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/koooyooo/soi-go/pkg/fileio"

	"github.com/koooyooo/soi-go/pkg/cli"
)

type (
	Repository interface {
		Store(context.Context, *cli.SoiVirtual) error
		StoreAll(context.Context, *cli.SoiVirtualBucket) error
		LoadAll(context.Context) (*cli.SoiVirtualBucket, error)
	}

	FileRepository struct {
		BasePath string
	}
)

func NewRepository() Repository {
	return FileRepository{
		BasePath: ".",
	}
}

func (f FileRepository) Store(ctx context.Context, s *cli.SoiVirtual) error {
	sb, err := f.LoadAll(ctx)
	if err != nil {
		return nil
	}
	sb.Sois = append(sb.Sois, s)
	return f.StoreAll(ctx, sb)
}

func (f FileRepository) StoreAll(ctx context.Context, sb *cli.SoiVirtualBucket) error {
	b, err := json.MarshalIndent(sb, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(f.BasePath, "repo.json"), b, 0600)
}

func (f FileRepository) LoadAll(ctx context.Context) (*cli.SoiVirtualBucket, error) {
	path := path.Join(f.BasePath, "repo.json")
	if !fileio.FileExists(path) {
		return &cli.SoiVirtualBucket{}, nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var soiBucket cli.SoiVirtualBucket
	if err := json.Unmarshal(b, &soiBucket); err != nil {
		return nil, err
	}
	return &soiBucket, nil
}
