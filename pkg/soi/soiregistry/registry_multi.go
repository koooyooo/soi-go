package soiregistry

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/koooyooo/soi-go/pkg/fileio"

	"github.com/koooyooo/soi-go/pkg/soi"
)

type multiRegistry struct {
}

func (mr multiRegistry) Register(s *soi.Soi) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	soisDir, err := fileio.SoisDirPath()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(soisDir, s.Name), b, 0600)
}

func (mr multiRegistry) FindByName(name string) ([]*soi.Soi, error) {
	soisDir, err := fileio.SoisDirPath()
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(soisDir, func(path string, info os.FileInfo, err error) error {

		return nil
	})
	return nil, err // TODO Implement
}
