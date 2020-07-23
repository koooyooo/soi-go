package soiprompt

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soi"
)

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

// filePathsToSuggests はファイルパスの配列を元に Suggestの配列を生成します
func filePathsToSuggests(soisDir string, files []string, word string) []prompt.Suggest {
	var s []prompt.Suggest
	for _, f := range files {
		s = append(s, prompt.Suggest{
			Text:        strings.TrimPrefix(f, soisDir+"/"),
			Description: "",
		})
	}
	return prompt.FilterContains(s, word, true)
}
