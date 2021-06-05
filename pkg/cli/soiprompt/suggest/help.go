package suggest

import (
	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/suggest/common"
)

func HelpCmd(d prompt.Document) []prompt.Suggest {
	return common.EmptySuggests
}
