package repo

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/koooyooo/soi-go/pkg/cli"
)

type (
	Repository interface {
		Store(context.Context, *cli.SoiWithPath) error
		StoreAll(context.Context, *cli.SoiBucket) error
		LoadAll(context.Context) (*cli.SoiBucket, error)
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

func (f FileRepository) Store(ctx context.Context, s *cli.SoiWithPath) error {
	sb, err := f.LoadAll(ctx)
	if err != nil {
		return nil
	}
	sb.Sois = append(sb.Sois, s)
	return f.StoreAll(ctx, sb)
}

func (f FileRepository) StoreAll(ctx context.Context, sb *cli.SoiBucket) error {
	b, err := json.MarshalIndent(sb, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(f.BasePath, "repo.json"), b, 0600)
}

func (f FileRepository) LoadAll(ctx context.Context) (*cli.SoiBucket, error) {
	b, err := ioutil.ReadFile(path.Join(f.BasePath, "repo.json"))
	if err != nil {
		return nil, err
	}
	var soiBucket cli.SoiBucket
	if err := json.Unmarshal(b, &soiBucket); err != nil {
		return nil, err
	}
	return &soiBucket, nil
}
