package main

import (
	prompt "github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt"
)

func main() {
	p := prompt.New(
		soiprompt.Executor,
		soiprompt.Completer,
		prompt.OptionTitle("soi input"),
		prompt.OptionPrefix("soi> "),
		prompt.OptionPrefixTextColor(prompt.Black),
		//prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionScrollbarThumbColor(prompt.Black),
		prompt.OptionScrollbarBGColor(prompt.DarkGray),
		prompt.OptionInputTextColor(prompt.Black),
		prompt.OptionPreviewSuggestionTextColor(prompt.DarkBlue),
		prompt.OptionMaxSuggestion(15),
	)
	p.Run()
}
