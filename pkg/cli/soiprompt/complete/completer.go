package complete

import (
	"soi-go/pkg/cli/config"
	"soi-go/pkg/cli/service"
	"soi-go/pkg/cli/soiprompt/cache"
	"strings"

	"github.com/c-bata/go-prompt"
	"soi-go/pkg/model"
)

var EmptySuggests []prompt.Suggest

func NewCompleter(conf *config.Config, service service.Service, ca *cache.Cache, b *model.BucketRef) *Completer {
	return &Completer{
		conf:      conf,
		service:   service,
		cache:     ca,
		BucketRef: b,
	}
}

type Completer struct {
	conf    *config.Config
	service service.Service
	cache   *cache.Cache
	*model.BucketRef
}

// Completer はSuggest候補を提示することで補完を実施します
func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	switch {
	case hasPrefixes(text, "add ", "a "):
		return c.addCmd(d)
	case hasPrefixes(text, "mv "):
		return c.mvCmd(d)
	case hasPrefixes(text, "rm "):
		return c.rmCmd(d)
	case hasPrefixes(text, "dig ", "d "):
		return c.digCmd(d)
	case hasPrefixes(text, "list ", "ls", "l "):
		return c.listCmd(d)
	case hasPrefixes(text, "edit ", "e "):
		return c.editCmd(d)
	case hasPrefixes(text, "cb ", "c "):
		return c.cbCmd(d)
	case hasPrefixes(text, "tag ", "t "):
		return c.tagCmd(d)
	case hasPrefixes(text, "help ", "h "):
		return c.helpCmd(d)
	case hasPrefixes(text, "size ", "s "):
		return nil
	default:
		s := []prompt.Suggest{
			{Text: "add", Description: "(a)dd url"},
			{Text: "dig", Description: "(d)ig urls"},
			{Text: "list", Description: "(l)ist and filter urls"},
			{Text: "cb", Description: "change bucket"},
			{Text: "tag", Description: "add tags (TODO)"},
			{Text: "mv", Description: "move file to dir"},
			{Text: "rm", Description: "remove file or dir"},
			//{Text: "tags", Description: "lists up all tags"}, TODO implements as "list -t"
			{Text: "pull", Description: "pull urls from repos"},
			{Text: "push", Description: "push urls to repos"},
			{Text: "help", Description: "(h)elp document"},
			{Text: "size", Description: "show (s)ize of soi data"},
			{Text: "quit", Description: "(q)uit model"},
			{Text: "version", Description: "show (v)ersion"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
}

// hasPrefixes は引数の文字列に接頭語が含まれているものがあるかを調査します
func hasPrefixes(in string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(in, p) {
			return true
		}
	}
	return false
}
