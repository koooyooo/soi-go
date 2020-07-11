package soiprompt

import (
	"fmt"
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
		return suggestAddCmd(d.GetWordBeforeCursor())
	case "open", "o":
		cmd := strings.TrimPrefix(textBC, "open ")
		cmd = strings.TrimPrefix(cmd, "o ")
		return suggestOpenCmd(cmd)
	case "list", "l":
		return suggestListCmd(d.GetWordBeforeCursor())
	default:
		s := []prompt.Suggest{
			{Text: "add", Description: "Add"},
			{Text: "list", Description: "List"},
			{Text: "l", Description: "List"},
			{Text: "tags", Description: "Tags"},
			{Text: "open", Description: "Open"},
			{Text: "o", Description: "Open"},
			{Text: "tag", Description: "Tag"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
	return []prompt.Suggest{}
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
		fmt.Println(dir) // TODO
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
