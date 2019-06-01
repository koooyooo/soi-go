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
