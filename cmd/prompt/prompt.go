package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
)

func completer(d prompt.Document) []prompt.Suggest {
	//fmt.Println("Cmp", d.TextBeforeCursor())
	s := []prompt.Suggest{
		{Text: "add", Description: "Add"},
		{Text: "list", Description: "List"},
		{Text: "tags", Description: "Tags"},
		{Text: "open", Description: "Open"},
		{Text: "tag", Description: "Tag"},
	}
	if strings.HasPrefix(d.TextBeforeCursor(), "add") {
		s = []prompt.Suggest{}
	}
	if strings.HasPrefix(d.TextBeforeCursor(), "open") {
		s = []prompt.Suggest{
			{"hello", "Hello"},
			{"world", "World"},
		}
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	fmt.Println("Please select table.")
	p := prompt.New(executor, completer, prompt.OptionTitle("soi input"), prompt.OptionPrefix("soi> "))
	p.Run()
}

func executor(in string) {
	if in == "exit" {
		os.Exit(0)
	}
}
