package registory

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	commons "github.com/koooyooo/soi-go/comons"
	"golang.org/x/xerrors"

	"github.com/koooyooo/soi-go/model"
)

type localRegistry struct{}

func (l localRegistry) Load() (*model.SoiCup, error) {
	s := model.SoiCup{}
	b, err := ioutil.ReadFile(commons.SoisFilePath)
	if err != nil {
		return nil, xerrors.Errorf("error in reading [%s] %v", commons.SoisFilePath, err)
	}
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, xerrors.Errorf("error in unmarshalling content [%s] %v", string(b), err)
	}
	return &s, nil
}

func (l localRegistry) Store(s model.SoiCup) error {
	b, err := json.Marshal(s)
	if err != nil {
		return xerrors.Errorf("error in marshalling content %v %v", s, err)
	}
	var prettyBuff bytes.Buffer
	err = json.Indent(&prettyBuff, b, "", "  ")
	if err != nil {
		return xerrors.Errorf("filed in indent json %v", err)
	}
	err = ioutil.WriteFile(commons.SoisFilePath, prettyBuff.Bytes(), 0666)
	if err != nil {
		return xerrors.Errorf("filed in writing file [%s] %v", commons.SoisFilePath, err)
	}
	return nil
}
