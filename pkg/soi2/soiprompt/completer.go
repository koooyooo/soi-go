package soiprompt

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/koooyooo/soi-go/pkg/fileio"

	"github.com/c-bata/go-prompt"
)

var (
	EmptySuggests = []prompt.Suggest{}
)

// completer はSuggest候補を提示することで補完を実施します
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
	case hasPrefixes(text, "dig ", "d "):
		return suggestDigCmd(d)
	case hasPrefixes(text, "list ", "l "):
		return suggestListCmd(d)
	default:
		s := []prompt.Suggest{
			{Text: "add", Description: "(a)dd url"},
			{Text: "open", Description: "(o)pen urls (complement path by args)"},
			{Text: "dig", Description: "(d)ig urls (complement path by -> key)"},
			{Text: "list", Description: "(l)ist urls and filter them"},
			{Text: "tag", Description: "add tags (not implemented yet)"},
			{Text: "mv", Description: "move file to dir"},
			{Text: "rm", Description: "remove file or dir"},
			//{Text: "tags", Description: "lists up all tags"}, TODO implements as "list -t"
			{Text: "quit", Description: "(q)uit soi"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
	return EmptySuggests
}

// suggestAddCmd はaddコマンド系のSuggestを提示します
func suggestAddCmd(d prompt.Document) []prompt.Suggest {
	// option探索
	if strings.HasPrefix(d.GetWordBeforeCursor(), "-") {
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
		soiRoot, err := fileio.SoisDirPath()
		if err != nil {
			log.Fatal(err)
		}
		dirs, err := listDirsRecursively(soiRoot, false)
		for _, d := range dirs {
			suggests = append(suggests, prompt.Suggest{
				Text:        strings.TrimPrefix(d, soiRoot+"/"),
				Description: "",
			})
		}
		return suggests
	}
	if strings.HasSuffix(d.Text, " ") {
		return []prompt.Suggest{
			{Text: "https://", Description: "target url"},
			{Text: "-n", Description: "name of the url"},
			{Text: "-d", Description: "dir to which soi store"},
		}
	}
	return EmptySuggests
}

// suggestMvCmd はmvコマンド系のSuggestを提示します
func suggestMvCmd(d prompt.Document) []prompt.Suggest {
	text := d.Text
	is2ndArg := 2 < len(strings.Split(text, " "))

	word := strings.TrimPrefix(d.GetWordBeforeCursor(), "mv ")

	dir, err := fileio.SoisDirPath()
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	if is2ndArg {
		files, err = listDirsRecursively(dir, true)
	} else {
		files, err = listFilesRecursively(dir)
	}
	if err != nil {
		log.Fatal(err)
	}
	return filePathsToSuggests(dir, files, word)
}

// suggestRmCmd はrmコマンド系のSuggestを提示します
func suggestRmCmd(d prompt.Document) []prompt.Suggest {
	word := strings.TrimPrefix(d.GetWordBeforeCursor(), "rm ")

	dir, err := fileio.SoisDirPath()
	if err != nil {
		log.Fatal(err)
	}
	var fileDirs []string
	// ファイル系を追加
	files, err := listFilesRecursively(dir)
	if err != nil {
		log.Fatal(err)
	}
	fileDirs = append(fileDirs, files...)
	// ディレクトリ系を追加
	dirs, err := listDirsRecursively(dir, false)
	if err != nil {
		log.Fatal(err)
	}
	fileDirs = append(fileDirs, dirs...)

	sort.Sort(sort.StringSlice(fileDirs))

	return filePathsToSuggests(dir, fileDirs, word)
}

// suggestOpenCmd はopenコマンド系のSuggestを提示します
func suggestOpenCmd(d prompt.Document) []prompt.Suggest {
	// option探索
	if strings.HasPrefix(d.GetWordBeforeCursor(), "-") {
		return []prompt.Suggest{
			{Text: "-f", Description: "open w/ firefox"},
			{Text: "-s", Description: "open w/ safari"},
		}
	}
	input := d.TextBeforeCursor()
	inputs := strings.Split(input, " ")

	if len(inputs) < 2 {
		return EmptySuggests
	}
	// 利用しないフラグもパースの関係上宣言を行う
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	flags.Bool("f", false, "open w/ firefox")
	flags.Bool("s", false, "open w/ safari")
	flags.Parse(inputs[1:])

	path := filepath.Join(flags.Args()...)

	// 相対パスを元にファイルを抽出
	soisDir, _ := fileio.SoisDirPath()
	return suggestByPath(soisDir, filepath.Join(soisDir, path), d.GetWordBeforeCursor(), false)

}

// suggestDigCmd はppコマンド系のSuggestを提示します
func suggestDigCmd(d prompt.Document) []prompt.Suggest {
	// option探索
	if strings.HasPrefix(d.GetWordBeforeCursor(), "-") {
		return []prompt.Suggest{
			{Text: "-f", Description: "dig w/ firefox"},
			{Text: "-s", Description: "dig w/ safari"},
		}
	}
	input := d.TextBeforeCursor()
	inputs := strings.Split(input, " ")

	flags := flag.NewFlagSet("dig", flag.PanicOnError)
	flags.Bool("f", false, "open w/ firefox")
	flags.Bool("s", false, "open w/ safari")
	flags.Parse(inputs[1:])

	soisDir, _ := fileio.SoisDirPath()
	return suggestByPath(soisDir, filepath.Join(soisDir, flags.Arg(0)), d.GetWordBeforeCursor(), true)
}

func suggestByPath(soisDir, path, input string, showDir bool) []prompt.Suggest {
	var found bool
	isDir, err := fileio.IsDir(strings.TrimSuffix(path, "/"))
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			found = false
			// 対象ファイルが見つからないだけの場合はスルー
		default:
			log.Fatal(err)
		}
	} else {
		found = true
	}
	switch {
	case !found || isDir:
		if !found {
			path = toLeafDirPath(path)
		}
		dirs, err := listFileDirs(path, showDir, true)
		if err != nil {
			log.Fatal(err)
		}
		return filePathsToSuggests(soisDir, dirs, input)
	}
	return EmptySuggests
}

// suggestListCmd はlistコマンド系のSuggestを提示します
func suggestListCmd(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest
	dir, _ := fileio.SoisDirPath()
	files, err := listFilesRecursively(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		s = append(s, prompt.Suggest{
			Text:        strings.TrimPrefix(f, dir+"/"),
			Description: "",
		})
	}
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}
