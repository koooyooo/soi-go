package soiprompt

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soi"
)

func hasPrefixes(in string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(in, p) {
			return true
		}
	}
	return false
}

func finalDirFromPath(path string) string {
	lastSlashIdx := strings.LastIndex(path, "/")
	return path[0:lastSlashIdx]
}

func listFileDirs(dir string, lastSlash bool) ([]string, error) {
	fileinfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, fi := range fileinfos {
		path := fi.Name()
		if fi.IsDir() && lastSlash {
			path = filepath.Join(dir, path) + "/"
		} else {
			path = filepath.Join(dir, path)
		}
		paths = append(paths, path)
	}
	return paths, nil
}

// listFilesRecursively はsoiRoot配下のファイルを再帰的に追加して Suggestを抽出します
func listFilesRecursively(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func listDirsRecursively(dir string, lastSlash bool) ([]string, error) {
	soiRoot, err := soi.SoisDirPath()
	if err != nil {
		return nil, err
	}
	var files []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != soiRoot {
			if lastSlash {
				path = path + "/"
			}
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
