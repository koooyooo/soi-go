package repo

import "github.com/koooyooo/soi-go/pkg/cli"

type (
	Repositry interface {
		Store(cli.SoiBucket) error
		Load(user, context string) (cli.SoiBucket, error)
	}
)
