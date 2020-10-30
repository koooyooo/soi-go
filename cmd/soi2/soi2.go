package main

import (
	prompt "github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soi2/soiprompt"
)

func main() {
	p := prompt.New(
		soiprompt.Executor,
		soiprompt.Completer,
		prompt.OptionTitle("soi input"),
		prompt.OptionPrefix("soi> "),
		prompt.OptionPrefixTextColor(prompt.Black),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionScrollbarBGColor(prompt.Blue),
		prompt.OptionInputTextColor(prompt.Black),
		prompt.OptionPreviewSuggestionTextColor(prompt.DarkBlue),
		prompt.OptionMaxSuggestion(15),
	)
	p.Run()
}
