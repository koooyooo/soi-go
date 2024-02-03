package complete

import (
	"fmt"
	"github.com/c-bata/go-prompt"
)

func (c *Completer) editCmd(d prompt.Document) []prompt.Suggest {
	fmt.Println("edit...") // TODO
	return c.baseList(d, "edit", "e")
}
