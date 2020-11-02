package repo

import "github.com/koooyooo/soi-go/pkg/cli"

type (
	Repository interface {
		Store(cli.SoiBucket) error
		Load(user, context string) (cli.SoiBucket, error)
	}
)
