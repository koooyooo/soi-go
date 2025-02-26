package complete

import (
	"github.com/c-bata/go-prompt"
)

func (c *Completer) helpCmd(_ prompt.Document) []prompt.Suggest {
	return EmptySuggests
}
