package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/koooyooo/soi-go/pkg/soi"

	"github.com/c-bata/go-prompt"
)

type SoiData struct {
	Name    string   `json:"name""`
	URI     string   `json:"uri"`
	Tags    []string `json:"tags"`
	Created string   `json"created"`
}

// PreviousTarget は事前の実行結果を記録します
type PreviousTarget struct {
	Path     string
	Suggests []prompt.Suggest
}

// filter は入力を元にSuggestをフィルタリングします
func (p PreviousTarget) filter(lastInput string) []prompt.Suggest {
	var filtered []prompt.Suggest
	for _, s := range p.Suggests {
		if strings.HasPrefix(s.Text, lastInput) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

var previousTarget PreviousTarget

func main() {
	fmt.Println("Please select table.")
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("soi input"),
		prompt.OptionPrefix("soi> "),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
	)
	p.Run()
}

func executor(in string) {
	fmt.Printf("EXEC: %s\n", in)
	cmd := strings.Split(in, " ")[0]
	subCmd := strings.TrimPrefix(in, cmd+" ")
	switch cmd {
	case "exit":
		os.Exit(0)
	case "open", "list":
		relPath := strings.ReplaceAll(subCmd, " ", "/")
		err := open(relPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func open(relPath string) error {
	dir, err := soi.SoisDirPath()
	if err != nil {
		return err
	}
	fullPath := dir + "/" + relPath
	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	var soi SoiData
	err = json.Unmarshal(b, &soi)
	if err != nil {
		return err
	}
	err = exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", soi.URI).Start()
	if err != nil {
		return err
	}
	return nil
}

// completer は補完を実施します
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
	cmd := strings.Split(textBC, " ")[0]
	switch cmd {
	case "add":
		return suggestAddCmd(d.GetWordBeforeCursor())
	case "open":
		input := strings.TrimPrefix(textBC, "open ")
		return suggestOpenCmd(input)
	case "list":
		return suggestListCmd(d.GetWordBeforeCursor())
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func suggestAddCmd(input string) []prompt.Suggest {
	return []prompt.Suggest{}
}

// suggestOpenCmd は指定した相対Path(soiRoot以降)を元に Suggestを抽出します
func suggestOpenCmd(input string) []prompt.Suggest {
	path := strings.ReplaceAll(input, " ", "/")
	// 相対パスを元にファイルを抽出
	rootDir, _ := soi.SoisDirPath()
	var dir string
	if path != " " {
		dir = rootDir + "/" + path
	}
	files, _ := ioutil.ReadDir(dir)

	// ファイルが存在しない場合は直前に保管したSuggestを提示 FIXME: Pathを戻したときの挙動に対応出来ていない
	if len(files) == 0 {
		idx := strings.LastIndex(path, "/")
		if idx != -1 {
			prevDir := path[0:idx]
			fmt.Println("PD", prevDir)
		}

		pathElm := strings.Split(path, "/")
		lastInput := pathElm[len(pathElm)-1]
		return previousTarget.filter(lastInput)
	}

	// ファイルが存在する場合は候補に保存の上、直前のSuggestとして保管
	var s []prompt.Suggest
	for _, f := range files {
		s = append(s, prompt.Suggest{Text: f.Name(), Description: ""})
	}
	previousTarget = PreviousTarget{
		Path:     dir,
		Suggests: s,
	}
	return s
}

// listFiles はsoiRoot配下のファイルを再帰的に追加して Suggestを抽出します
func listFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, listFiles(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}

func suggestListCmd(input string) []prompt.Suggest {
	s := []prompt.Suggest{}
	dir, _ := soi.SoisDirPath()
	files := listFiles(dir)
	for _, f := range files {
		s = append(s, prompt.Suggest{
			Text:        strings.TrimPrefix(f, dir+"/"),
			Description: "",
		})
	}
	return prompt.FilterContains(s, input, true)
}
