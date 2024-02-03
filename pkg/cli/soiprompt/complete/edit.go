package complete

import (
	"github.com/c-bata/go-prompt"
)

func (c *Completer) editCmd(d prompt.Document) []prompt.Suggest {
	return c.baseList(d, "edit", "e")
}
