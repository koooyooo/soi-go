package execute

import (
	"fmt"
	"strings"

	"github.com/koooyooo/soi-go/pkg/cli/config"
	"github.com/koooyooo/soi-go/pkg/cli/service"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/execute/registry"
	"github.com/koooyooo/soi-go/pkg/model"
)

func NewExecutor(c *config.Config, s service.Service, r *model.BucketRef) *Executor {
	return &Executor{
		Conf:      c,
		Service:   s,
		BucketRef: r,
	}
}

type Executor struct {
	Conf    *config.Config
	Service service.Service
	*model.BucketRef
}

// Executor は入力されたコマンドに応じた処理を行います
func (e *Executor) Execute(in string) {
	in = strings.Trim(in, " ")
	cmd := strings.Split(in, " ")[0]
	switch cmd {
	case "add", "a":
		op(e.add, in)
	case "mv":
		op(e.mv, in)
	case "rm":
		op(e.rm, in)
	case "cb":
		op(e.cb, in)
	case "tag", "t":
		op(e.tag, in)
	case "open", "o", "list", "ls", "l", "dig", "d":
		op(e.open, in)
	case "help", "h":
		op(e.help, in)
	case "pull":
		if err := registry.Pull(e.Conf, e.BucketRef.Bucket, in); err != nil {
			fmt.Println(err)
			return
		}
	case "push":
		if err := registry.Push(e.Conf, e.BucketRef.Bucket, in); err != nil {
			fmt.Println(err)
			return
		}
	case "quit", "q", "exit":
		op(e.quit, in)
	case "version", "v":
		op(e.version, in)
	}
}

func op(f func(string) error, in string) {
	if err := f(in); err != nil {
		fmt.Println(err)
		return
	}
}
