package execute

import (
	"fmt"
	"strings"

	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/execute/registry"
)

// Executor は入力されたコマンドに応じた処理を行います
func Executor(in string) {
	in = strings.Trim(in, " ")
	cmd := strings.Split(in, " ")[0]
	switch cmd {
	case "add", "a":
		if err := add(in); err != nil {
			fmt.Println(err)
			return
		}
	case "mv":
		if err := mv(in); err != nil {
			fmt.Println(err)
			return
		}
	case "rm":
		if err := rm(in); err != nil {
			fmt.Println(err)
			return
		}
	case "cb":
		if err := cb(in); err != nil {
			fmt.Println(err)
			return
		}
	case "open", "o", "list", "ls", "l", "dig", "d":
		if err := open(in); err != nil {
			fmt.Println(err)
			return
		}
	case "help", "h":
		if err := help(in); err != nil {
			fmt.Println(err)
			return
		}
	case "pull":
		if err := registry.Pull(in); err != nil {
			fmt.Println(err)
			return
		}
	case "push":
		if err := registry.Push(in); err != nil {
			fmt.Println(err)
			return
		}
	case "quit", "q", "exit":
		if err := quit(in); err != nil {
			fmt.Println(err)
			return
		}
	}
}
