package soiprompt

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soi"
)

// completer は補完を実施します
func Completer(d prompt.Document) []prompt.Suggest {
	textBC := d.TextBeforeCursor()
	cmd := strings.Split(textBC, " ")[0]
	switch cmd {
	case "add":
		return suggestAddCmd(d)
	case "open", "o":
		return suggestOpenCmd(d)
	case "list", "l":
		return suggestListCmd(d.GetWordBeforeCursor())
	default:
		s := []prompt.Suggest{
			{Text: "add", Description: "adds url"},
			{Text: "list", Description: "(l)ists urls and filter them"},
			{Text: "l", Description: "(l)ists urls and filter them"},
			{Text: "tags", Description: "lists up all tags"},
			{Text: "open", Description: "(0)pens specified url"},
			{Text: "o", Description: "(0)pens specified url"},
			{Text: "tag", Description: "adds tags"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
	return []prompt.Suggest{}
}

// suggestAddCmd は
func suggestAddCmd(d prompt.Document) []prompt.Suggest {
	if d.GetWordBeforeCursor() == "-" {
		return []prompt.Suggest{
			{Text: "-n", Description: "logical name of the url"},
		}
	}
	return []prompt.Suggest{}
}

// suggestOpenCmd は指定した相対Path(soiRoot以降)を元に Suggestを抽出します
func suggestOpenCmd(d prompt.Document) []prompt.Suggest {
	input := d.TextBeforeCursor()
	input = strings.TrimPrefix(input, "open ")
	input = strings.TrimPrefix(input, "o ")

	path := strings.ReplaceAll(input, " ", "/")
	// 相対パスを元にファイルを抽出
	rootDir, _ := soi.SoisDirPath()
	var dir string
	if path != " " {
		dir = rootDir + "/" + path
	}
	files, _ := ioutil.ReadDir(dir)

	if strings.HasSuffix(strings.Trim(input, " "), ".json") {
		return []prompt.Suggest{}
	}

	// 絞り込み中も候補を表示する処理
	if len(files) == 0 {
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
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			files, err := listFiles(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			paths = append(paths, files...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths, nil
}
