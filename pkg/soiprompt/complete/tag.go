package complete

import (
	"github.com/c-bata/go-prompt"
)

// tagCmd は
func (c *Completer) tagCmd(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}
