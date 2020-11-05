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
		Store(context.Context, *cli.SoiBucket) error
		Load(context.Context) (*cli.SoiBucket, error)
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

func (f FileRepository) Store(ctx context.Context, sb *cli.SoiBucket) error {
	b, err := json.Marshal(sb)
	if err != nil {
		return err
	}
	ioutil.WriteFile(path.Join(f.BasePath, "repo.json"), b, 0600)
	return nil
}

func (f FileRepository) Load(ctx context.Context) (*cli.SoiBucket, error) {
	b, err := ioutil.ReadFile(path.Join(f.BasePath, "repo.json"))
	if err != nil {
		return nil, err
	}
	var soiBucket *cli.SoiBucket
	if err := json.Unmarshal(b, soiBucket); err != nil {
		return nil, err
	}
	return soiBucket, nil
}
