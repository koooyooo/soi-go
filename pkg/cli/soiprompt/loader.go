package soiprompt

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"

	"github.com/koooyooo/soi-go/pkg/soi"
)

func readSoiData(filepath string) (*soi.SoiData, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var sd soi.SoiData
	if err := json.Unmarshal(b, &sd); err != nil {
		return nil, err
	}
	return &sd, nil
}

func loadSoiWithPath(files []string) ([]*soi.SoiWithPath, error) {
	var wg sync.WaitGroup
	wg.Add(len(files))

	var ss = make([]*soi.SoiWithPath, len(files))
	for i, f := range files {
		go func(idx int, fp string) {
			sd, err := readSoiData(fp)
			if err != nil {
				log.Fatalf("failed in load sd: %s", err.Error())
			}
			ss[idx] = &soi.SoiWithPath{
				SoiData: sd,
				Path:    fp,
			}
			wg.Done()
		}(i, f)
	}
	wg.Wait()
	return ss, nil
}
