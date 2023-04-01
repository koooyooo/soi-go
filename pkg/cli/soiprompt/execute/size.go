package execute

import (
	"fmt"
	"golang.org/x/net/context"
)

func (e *Executor) size(_ string) error {
	size, err := e.Service.Size(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("size: %d\n", size)
	return nil
}
