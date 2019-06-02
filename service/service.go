package service

import (
	"github.com/koooyooo/soi-go/model"
	"github.com/koooyooo/soi-go/registory"
)

func Add(name, uri string, tags []string) (*model.Soi, error) {
	//fmt.Printf("soi add name=%s, url=%s \n", name, uri)
	sois, err := registory.Load()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return &soi, nil
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
	targetSois := model.FilterByName(sois.Sois, name)
	if 0 == len(targetSois) {
		return nil, false, nil
	}
	return &targetSois[0], true, nil
}

func Remove(name string) error {
	sois, err := registory.Load()
	if err != nil {
		return err
	}
	removedSois := model.FilterByExcludeName(sois.Sois, name)
	sois.Sois = removedSois
	err = registory.Store(*sois)
	if err != nil {
		return err
	}
	return nil
}

func Tag(name string, tags []string) (*model.Soi, bool, error) {
	target, ok, err := Get(name)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	soi, err := Add(target.Name, target.Uri, tags)
	if err != nil {
		return nil, true, err
	}
	return soi, true, nil
}
