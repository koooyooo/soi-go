package main

import (
	"github.com/koooyooo/soi-go/pkg/soi2/soiprompt"

	"github.com/c-bata/go-prompt"
)

func main() {
	p := prompt.New(
		soiprompt.Executor,
		soiprompt.Completer,
		prompt.OptionTitle("soi input"),
		prompt.OptionPrefix("soi> "),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionScrollbarBGColor(prompt.Blue),
		prompt.OptionMaxSuggestion(15),
	)
	p.Run()
}
