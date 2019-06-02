package service

import (
	"github.com/koooyooo/soi-go/model"
	"github.com/koooyooo/soi-go/registory"
)

func Add(name, uri string, tags []string) (*model.Soi, error) {
	soiCup, err := registory.Load()
	if err != nil {
		return nil, err
	}
	soi := model.Soi{
		Name: name,
		Uri:  uri,
		Tags: tags,
	}
	newSois := model.FilterByExcludeName(soiCup.Sois, name)
	newSois = append(newSois, soi)
	soiCup.Sois = newSois
	err = registory.Store(*soiCup)
	if err != nil {
		return nil, err
	}
	return &soi, nil
}

func Search(namepart string) ([]model.Soi, error) {
	soiCup, err := registory.Load()
	if err != nil {
		return nil, err
	}
	if namepart == "" {
		return soiCup.Sois, nil
	}
	result := model.FilterByNamePart(soiCup.Sois, namepart)
	return result, nil
}

func Get(name string) (*model.Soi, bool, error) {
	soiCup, err := registory.Load()
	if err != nil {
		return nil, false, err
	}
	targetSois := model.FilterByName(soiCup.Sois, name)
	if 0 == len(targetSois) {
		return nil, false, nil
	}
	return &targetSois[0], true, nil
}

func Remove(name string) error {
	soiCup, err := registory.Load()
	if err != nil {
		return err
	}
	removedSois := model.FilterByExcludeName(soiCup.Sois, name)
	soiCup.Sois = removedSois
	err = registory.Store(*soiCup)
	if err != nil {
		return err
	}
	return nil
}

func Tag(name string, tags []string) (*model.Soi, bool, error) {
	soiBefore, ok, err := Get(name)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	soi, err := Add(soiBefore.Name, soiBefore.Uri, tags)
	if err != nil {
		return nil, true, err
	}
	return soi, true, nil
}
