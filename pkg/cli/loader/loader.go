package loader

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"github.com/koooyooo/soi-go/pkg/cli/constant"
	"strings"
	"sync"

	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"

	"github.com/koooyooo/soi-go/pkg/common/file"

	"github.com/koooyooo/soi-go/pkg/common/hash"

	"github.com/koooyooo/soi-go/pkg/model"
)

var isSoiFile = func(soiPath string) bool {
	return strings.HasSuffix(soiPath, ".json")
}

func LoadSois(filepath string) ([]*model.SoiData, error) {
	files, err := utils.ListFilesRecursively(filepath)
	if err != nil {
		log.Fatal(err)
	}
	sois, err := loadFilteredSoiDataArray(files)
	if err != nil {
		log.Fatal(err)
	}
	return sois, nil
}

// loadFilteredSoiDataArray は指定されたファイルパスの配列からSoiData(WithPath)の配列をロードします
func loadFilteredSoiDataArray(files []string) ([]*model.SoiData, error) {
	var filtered []string
	for _, f := range files {
		if !isSoiFile(f) {
			fmt.Printf("[Warn] found unknown format file: %s", f)
			continue
		}
		filtered = append(filtered, f)
	}
	return loadSoiDataArray(filtered)
}

func loadSoiDataArray(files []string) ([]*model.SoiData, error) {
	var wg sync.WaitGroup
	wg.Add(len(files))

	var ss = make([]*model.SoiData, len(files))
	for i, f := range files {
		go func(idx int, fp string) {
			defer wg.Done()
			sd, err := LoadSoiData(fp)
			if err != nil {
				fmt.Printf("failed in load sd: %s", err.Error())
				return
			}
			ss[idx] = sd
		}(i, f)
	}
	wg.Wait()
	return ss, nil
}

// loadSoiData は指定されたファイルパスよりSoiデータをロードします
func LoadSoiData(filepath string) (*model.SoiData, error) {
	filepath = addJSONSuffix(filepath)
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var sd model.SoiData
	if err := json.Unmarshal(b, &sd); err != nil {
		return nil, err
	}
	// complement fields // TODO fix this
	soisDir, err := constant.SoisDir()
	if err != nil {
		return nil, err
	}
	path := filepath
	path = strings.TrimPrefix(path, soisDir+"/")
	idxBucketTail := strings.Index(path, "/")
	path = path[idxBucketTail+1:]
	path = strings.TrimSuffix(sd.Path, ".json")
	sd.Path = path
	if sd.Hash == "" {
		sd.Hash, err = hash.Sha1(sd.URI)
		if err != nil {
			return nil, err
		}
	}
	return &sd, nil
}

func StoreSoiData(filepath string, s *model.SoiData) error {
	filepath = addJSONSuffix(filepath)
	ub, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath, ub, 0600); err != nil {
		return err
	}
	return nil
}

func Exists(filepath string) bool {
	filepath = addJSONSuffix(filepath)
	return file.Exists(filepath)
}

// 末尾に ".json" を追加します
func addJSONSuffix(path string) string {
	if !strings.HasSuffix(path, ".json") {
		path = path + ".json"
	}
	return path
}
