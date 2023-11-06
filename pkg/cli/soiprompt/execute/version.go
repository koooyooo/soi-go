package execute

import (
	"fmt"

	"soi-go/pkg/cli/constant"
)

func (e *Executor) version(_ string) error {
	fmt.Printf("version: %s \n", constant.Version)
	return nil
}
