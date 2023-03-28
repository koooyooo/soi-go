package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/model"
)

func ListFileDirs(dir string, addDir, dirLastSlash bool) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, fi := range fileInfos {
		path := fi.Name()
		if fi.IsDir() && dirLastSlash {
			if addDir {
				path = filepath.Join(dir, path) + "/"
			} else {
				path = path + "/"
			}
		} else {
			if addDir {
				path = filepath.Join(dir, path)
			}
		}
		paths = append(paths, path)
	}
	return paths, nil
}

// ListFilesRecursively は引数のdir配下を再帰的に走査してFileのPathを収集します
func ListFilesRecursively(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// ListDirsRecursively は引数のdir配下を再帰的に走査してDirectoryのPathを収集します
func ListDirsRecursively(dir string, b *model.Bucket, lastSlash bool) ([]string, error) {
	soiRoot, err := b.Path()
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

// FilePathsToSuggests は引数のfilePath配列を元に Suggestの配列を生成します
func FilePathsToSuggests(soisDir string, filePaths []string, word string) []prompt.Suggest {
	var s []prompt.Suggest
	for _, f := range filePaths {
		s = append(s, prompt.Suggest{
			Text:        strings.TrimPrefix(f, soisDir+"/"),
			Description: "",
		})
	}
	return prompt.FilterContains(s, word, true)
}

func FilePathsToSuggestsNoEx(soisDir string, filePaths []string, word string) []prompt.Suggest {
	suggests := FilePathsToSuggests(soisDir, filePaths, word)
	var suggestsNoEx []prompt.Suggest
	for _, s := range suggests {
		s.Text = strings.TrimSuffix(s.Text, ".json")
		suggestsNoEx = append(suggestsNoEx, s)
	}
	return suggestsNoEx
}
