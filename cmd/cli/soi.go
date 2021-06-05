package main

import (
	"log"
	"os"

	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/execute"

	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/suggest"

	prompt "github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/cli/constant"
	"github.com/koooyooo/soi-go/pkg/fileio"
)

func main() {
	soisDir, err := constant.LocalBucket.Path()
	if err != nil {
		log.Fatal(err)
	}
	if !fileio.Exists(soisDir) {
		if err := os.MkdirAll(soisDir, 0700); err != nil {
			log.Fatal(err)
		}
	}

	p := prompt.New(
		execute.Executor,
		suggest.Completer,
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
