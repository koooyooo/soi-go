package soiprompt

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soi"
)

var (
	EmptySuggests = []prompt.Suggest{}
)

// completer は補完を実施します
func Completer(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	switch {
	case hasPrefixes(text, "add ", "a "):
		return suggestAddCmd(d)
	case hasPrefixes(text, "mv "):
		return suggestMvCmd(d)
	case hasPrefixes(text, "open ", "o "):
		return suggestOpenCmd(d)
	case hasPrefixes(text, "list ", "l "):
		return suggestListCmd(d.GetWordBeforeCursor())
	default:
		s := []prompt.Suggest{
			{Text: "add", Description: "(a)dds url"},
			{Text: "mv", Description: "move file to dir"},
			{Text: "list", Description: "(l)ists urls and filter them"},
			{Text: "tags", Description: "lists up all tags"},
			{Text: "open", Description: "(o)pens specified url"},
			{Text: "tag", Description: "adds tags"},
			{Text: "quit", Description: "(q)uit soi"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
	return EmptySuggests
}

func hasPrefixes(in string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(in, p) {
			return true
		}
	}
	return false
}

// suggestAddCmd
func suggestAddCmd(d prompt.Document) []prompt.Suggest {
	// option探索
	if d.GetWordBeforeCursor() == "-" {
		return []prompt.Suggest{
			{Text: "-n", Description: "name of the url"},
			{Text: "-d", Description: "dir to which soi store"},
		}
	}
	if strings.HasSuffix(d.Text, "-n ") {
		return EmptySuggests
	}
	// dir探索
	if strings.HasSuffix(d.Text, "-d ") {
		var suggests []prompt.Suggest
		soiRoot, err := soi.SoisDirPath()
		if err != nil {
			log.Fatal(err)
		}
		dirs, err := listDirs(soiRoot)
		for _, d := range dirs {
			suggests = append(suggests, prompt.Suggest{
				Text:        strings.TrimPrefix(d, soiRoot+"/"),
				Description: "",
			})
		}
		return suggests
	}
	if strings.HasSuffix(d.Text, " ") {
		return []prompt.Suggest{{
			Text:        "https://",
			Description: "",
		}}
	}
	return EmptySuggests
}

// suggestMvCmd
func suggestMvCmd(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest

	text := d.Text
	is2ndArg := 2 < len(strings.Split(text, " "))

	word := d.GetWordBeforeCursor()
	word = strings.TrimPrefix(word, "mv ")

	dir, _ := soi.SoisDirPath()
	var files []string
	var err error
	if is2ndArg {
		files, err = listDirs(dir)
	} else {
		files, err = listFiles(dir)
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		s = append(s, prompt.Suggest{
			Text:        strings.TrimPrefix(f, dir+"/"),
			Description: "",
		})
	}
	return prompt.FilterContains(s, word, true)

	return EmptySuggests
}

// suggestOpenCmd は指定した相対Path(soiRoot以降)を元に Suggestを抽出します
func suggestOpenCmd(d prompt.Document) []prompt.Suggest {
	input := d.TextBeforeCursor()
	// コマンド部分を除去
	input = strings.TrimPrefix(input, "open ")
	input = strings.TrimPrefix(input, "o ")

	path := strings.ReplaceAll(input, " ", "/")
	// 相対パスを元にファイルを抽出
	rootDir, _ := soi.SoisDirPath()
	var dir string
	if path != " " {
		dir = rootDir + "/" + path
	}
	files, err := ioutil.ReadDir(dir)
	// 絞り込み中も候補を表示する処理
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		pathElm := strings.Split(path, "/")
		lastInput := pathElm[len(pathElm)-1]
		return previousTarget.filter(lastInput)
	}

	if strings.HasSuffix(strings.Trim(input, " "), ".json") {
		return []prompt.Suggest{}
	}

	if len(files) == 0 {
		return EmptySuggests
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

// suggestListCmd は"list"コマンドの制御を行います
func suggestListCmd(input string) []prompt.Suggest {
	var s []prompt.Suggest
	dir, _ := soi.SoisDirPath()
	files, err := listFiles(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		s = append(s, prompt.Suggest{
			Text:        strings.TrimPrefix(f, dir+"/"),
			Description: "",
		})
	}
	return prompt.FilterContains(s, input, true)
}

// listFiles はsoiRoot配下のファイルを再帰的に追加して Suggestを抽出します
func listFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func listDirs(dir string) ([]string, error) {
	soiRoot, err := soi.SoisDirPath()
	if err != nil {
		return nil, err
	}
	var files []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != soiRoot {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
