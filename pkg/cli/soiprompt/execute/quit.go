package execute

import (
	"fmt"
	"os"
)

func quit(in string) error {
	fmt.Println("bye!")
	os.Exit(0)
	return nil
}
