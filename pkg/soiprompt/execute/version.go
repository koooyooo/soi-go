package execute

import (
	"fmt"

	"github.com/koooyooo/soi-go/pkg/constant"
)

func (e *Executor) version(_ string) error {
	fmt.Printf("version: %s \n", constant.Version)
	return nil
}
