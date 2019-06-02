package registory

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/koooyooo/soi-go/model"
)

func Load() (*model.SoiCup, error) {
	s := model.SoiCup{}
	b, err := ioutil.ReadFile("sois.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func Store(s model.SoiCup) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	var prettyBuff bytes.Buffer
	err = json.Indent(&prettyBuff, b, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("sois.json", prettyBuff.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}
