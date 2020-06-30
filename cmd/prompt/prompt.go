package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/koooyooo/soi-go/pkg/soi"

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
	textBC := d.TextBeforeCursor()
	if strings.HasPrefix(textBC, "add") {
		s = []prompt.Suggest{}
	}
	if strings.HasPrefix(textBC, "open") {
		spaceSepPath := strings.TrimPrefix(textBC, "open ")
		path := strings.ReplaceAll(spaceSepPath, " ", "/")
		return readFileInfo(path)
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

type PreviousTarget struct {
	Path    string
	Suggest []prompt.Suggest
}

func (p PreviousTarget) filter(lastInput string) []prompt.Suggest {
	var filtered []prompt.Suggest
	for _, s := range p.Suggest {
		if strings.HasPrefix(s.Text, lastInput) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

var previousTarget PreviousTarget

func readFileInfo(path string) []prompt.Suggest {
	var s []prompt.Suggest
	dir, _ := soi.SoisDirPath()
	if path != "" && path != " " {
		dir = dir + "/" + path
	}
	files, _ := ioutil.ReadDir(dir)
	if len(files) != 0 {
		for _, f := range files {
			s = append(s, prompt.Suggest{Text: f.Name(), Description: ""})
		}
		previousTarget = PreviousTarget{
			Path:    dir,
			Suggest: s,
		}
	} else {
		pathElm := strings.Split(path, "/")
		lastInput := pathElm[len(pathElm)-1]
		s = previousTarget.filter(lastInput)
	}
	return s
}

func executor(in string) {
	fmt.Printf("EXEC: %s\n", in)
	if in == "exit" {
		os.Exit(0)
	}
}

func main() {
	fmt.Println("Please select table.")
	p := prompt.New(executor, completer, prompt.OptionTitle("soi input"), prompt.OptionPrefix("soi> "))
	p.Run()
}
