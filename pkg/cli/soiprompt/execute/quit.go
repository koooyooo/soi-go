package execute

import (
	"fmt"
	"os"
)

func (e *Executor) quit(_ string) error {
	fmt.Println("bye!")
	os.Exit(0)
	return nil
}
