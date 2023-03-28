package execute

import (
	"fmt"

	"github.com/koooyooo/soi-go/pkg/cli/constant"
)

func (e *Executor) version(_ string) error {
	fmt.Printf("version: %s \n", constant.Version)
	return nil
}
