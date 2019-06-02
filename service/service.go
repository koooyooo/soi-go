package service

import (
	"fmt"

	"github.com/koooyooo/soi-go/model"
	"github.com/koooyooo/soi-go/registory"
)

func Add(name, uri string, tags []string) error {
	//fmt.Printf("soi add name=%s, url=%s \n", name, uri)
	sois, err := registory.Load()
	if err != nil {
		return err
	}
	if sois.Contains(name) {
		sois.Remove(name)
	}
	soi := model.Soi{
		Name: name,
		Uri:  uri,
		Tags: tags,
	}
	sois.Add(soi)
	err = registory.Store(*sois)
	if err != nil {
		return err
	}
	fmt.Printf("stored: %v\n", soi)
	return nil
}

func Search(namepart string) ([]model.Soi, error) {
	sois, err := registory.Load()
	if err != nil {
		return nil, err
	}
	if namepart == "" {
		return sois.Sois, nil
	}
	result := model.FilterByNamePart(sois.Sois, namepart)
	return result, nil
}

func Get(name string) (*model.Soi, bool, error) {
	sois, err := registory.Load()
	if err != nil {
		return nil, false, err
	}
	for _, v := range sois.Sois {
		if v.Name == name {
			return &v, true, nil
		}
	}
	return nil, false, nil
}

func Remove(name string) error {
	sois, err := registory.Load()
	if err != nil {
		return err
	}
	var newSoiList []model.Soi
	for _, v := range sois.Sois {
		if v.Name != name {
			newSoiList = append(newSoiList, v)
		}
	}
	sois.Sois = newSoiList
	err = registory.Store(*sois)
	if err != nil {
		return err
	}
	return nil
}
