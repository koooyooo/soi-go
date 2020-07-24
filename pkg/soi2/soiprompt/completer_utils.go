package soiprompt

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/koooyooo/soi-go/pkg/soi"
)

// hasPrefixes は引数の文字列に接頭語が含まれているものがあるかを調査します
func hasPrefixes(in string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(in, p) {
			return true
		}
	}
	return false
}

// toLeafDirPath はPathを末端ディレクトリのPathに変換します
func toLeafDirPath(path string) string {
	lastSlashIdx := strings.LastIndex(path, "/")
	return path[0:lastSlashIdx]
}

// toStorablePath はPathを保存可能なPathに変換します
func toStorablePath(path string) string {
	// path として意味を持ってしまう文字列を変換します
	path = strings.ReplaceAll(path, " ", "_")
	path = strings.ReplaceAll(path, "/", "／")
	if !strings.HasSuffix(path, ".json") {
		path = path + ".json"
	}
	return path
}

// listFileDirs は
func listFileDirs(dir string, lastSlash bool) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, fi := range fileInfos {
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

// listFilesRecursively は引数のdir配下を再帰的に走査してFileのPathを収集します
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

// listDirsRecursively は引数のdir配下を再帰的に走査してDirectoryのPathを収集します
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

// filePathsToSuggests は引数のfilePath配列を元に Suggestの配列を生成します
func filePathsToSuggests(soisDir string, filePaths []string, word string) []prompt.Suggest {
	var s []prompt.Suggest
	for _, f := range filePaths {
		s = append(s, prompt.Suggest{
			Text:        strings.TrimPrefix(f, soisDir+"/"),
			Description: "",
		})
	}
	return prompt.FilterContains(s, word, true)
}
