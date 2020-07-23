package soiprompt

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
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
	case hasPrefixes(text, "rm "):
		return suggestRmCmd(d)
	case hasPrefixes(text, "open ", "o "):
		return suggestOpenCmd(d)
	case hasPrefixes(text, "pp"):
		return suggestPpCmd(d)
	case hasPrefixes(text, "list ", "l "):
		return suggestListCmd(d.GetWordBeforeCursor())
	default:
		s := []prompt.Suggest{
			{Text: "add", Description: "(a)dds url"},
			{Text: "list", Description: "(l)ists urls and filter them"},
			{Text: "mv", Description: "move file to dir"},
			{Text: "rm", Description: "remove file or dir"},
			//{Text: "tags", Description: "lists up all tags"}, TODO implements as "list -t"
			{Text: "open", Description: "(o)pens specified url"},
			{Text: "quit", Description: "(q)uit soi"},
			{Text: "tag", Description: "adds tags"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
	return EmptySuggests
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
		dirs, err := listDirs(soiRoot, false)
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
	text := d.Text
	is2ndArg := 2 < len(strings.Split(text, " "))

	word := strings.TrimPrefix(d.GetWordBeforeCursor(), "mv ")

	dir, err := soi.SoisDirPath()
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	if is2ndArg {
		files, err = listDirs(dir, true)
	} else {
		files, err = listFiles(dir)
	}
	if err != nil {
		log.Fatal(err)
	}
	return filePathsToSuggests(dir, files, word)
}

// suggestRmCmd
func suggestRmCmd(d prompt.Document) []prompt.Suggest {
	word := strings.TrimPrefix(d.GetWordBeforeCursor(), "rm ")

	dir, err := soi.SoisDirPath()
	if err != nil {
		log.Fatal(err)
	}
	var fileDirs []string
	// ファイル系を追加
	files, err := listFiles(dir)
	if err != nil {
		log.Fatal(err)
	}
	fileDirs = append(fileDirs, files...)
	// ディレクトリ系を追加
	dirs, err := listDirs(dir, false)
	if err != nil {
		log.Fatal(err)
	}
	fileDirs = append(fileDirs, dirs...)

	sort.Sort(sort.StringSlice(fileDirs))

	return filePathsToSuggests(dir, fileDirs, word)
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

func suggestPpCmd(d prompt.Document) []prompt.Suggest {
	input := d.TextBeforeCursor()
	input = strings.TrimPrefix(input, "pp ")

	rootDir, _ := soi.SoisDirPath()

	// case-1: 直近のフォルダ・ファイルに前方一致 => ファイル・フォルダの一覧を表示 (フォルダは末尾に "/"を付与)
	// case-2: 直近のフォルダに完全一致
	// case-3: 直近のファイルに完全一致

	fmt.Println(input, rootDir) // TODO Remote Later
	return EmptySuggests
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
