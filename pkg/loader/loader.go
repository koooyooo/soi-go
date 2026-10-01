package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/koooyooo/soi-go/pkg/common/file"
	"github.com/koooyooo/soi-go/pkg/common/hash"
	"github.com/koooyooo/soi-go/pkg/model"
	"github.com/koooyooo/soi-go/pkg/soiprompt/utils"
)

var isSoiFile = func(soiPath string) bool {
	return strings.HasSuffix(soiPath, ".json")
}

func LoadSois(bucketRoot string) ([]*model.SoiData, error) {
	files, err := utils.ListFilesRecursively(bucketRoot)
	if err != nil {
		return nil, err
	}
	return loadFilteredSoiDataArray(bucketRoot, files)
}

// loadFilteredSoiDataArray は指定されたファイルパスの配列からSoiDataの配列をロードします
func loadFilteredSoiDataArray(bucketRoot string, files []string) ([]*model.SoiData, error) {
	var filtered []string
	for _, f := range files {
		if !isSoiFile(f) {
			fmt.Printf("[Warn] found unknown format file: %s\n", f)
			continue
		}
		filtered = append(filtered, f)
	}
	return loadSoiDataArray(bucketRoot, filtered)
}

func loadSoiDataArray(bucketRoot string, files []string) ([]*model.SoiData, error) {
	var wg sync.WaitGroup
	wg.Add(len(files))

	var ss = make([]*model.SoiData, len(files))
	var firstErr error
	var mu sync.Mutex
	for i, f := range files {
		go func(idx int, fp string) {
			defer wg.Done()
			sd, err := LoadSoiData(bucketRoot, fp)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				fmt.Printf("failed in load sd: %s\n", err.Error())
				return
			}
			ss[idx] = sd
		}(i, f)
	}
	wg.Wait()
	if firstErr != nil {
		// 部分的に読めている場合でも、呼び出し側が扱えるよう成功分だけ返す
		var compact []*model.SoiData
		for _, s := range ss {
			if s != nil {
				compact = append(compact, s)
			}
		}
		return compact, nil
	}
	return ss, nil
}

// LoadSoiData は指定されたファイルパスよりSoiデータをロードします。
// Path は bucketRoot からの相対ディレクトリに正規化します。
func LoadSoiData(bucketRoot, filePath string) (*model.SoiData, error) {
	filePath = addJSONSuffix(filePath)
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var sd model.SoiData
	if err := json.Unmarshal(b, &sd); err != nil {
		return nil, err
	}

	rel, err := filepath.Rel(bucketRoot, filePath)
	if err != nil {
		return nil, fmt.Errorf("resolve relative path: %w", err)
	}
	rel = filepath.ToSlash(rel)
	dir := filepath.ToSlash(filepath.Dir(rel))
	if dir == "." {
		dir = ""
	}
	sd.Path = dir

	base := filepath.Base(rel)
	fileName := strings.TrimSuffix(base, ".json")
	if sd.Name == "" {
		sd.Name = fileName
	}

	if sd.Hash == "" {
		sd.Hash, err = hash.Sha1(sd.URI)
		if err != nil {
			return nil, err
		}
	}
	return &sd, nil
}

func StoreSoiData(filePath string, s *model.SoiData) error {
	filePath = addJSONSuffix(filePath)
	if err := os.MkdirAll(filepath.Dir(filePath), 0700); err != nil {
		return err
	}
	ub, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, ub, 0600)
}

func Exists(filePath string) bool {
	filePath = addJSONSuffix(filePath)
	return file.Exists(filePath)
}

// 末尾に ".json" を追加します
func addJSONSuffix(path string) string {
	if !strings.HasSuffix(path, ".json") {
		path = path + ".json"
	}
	return path
}
