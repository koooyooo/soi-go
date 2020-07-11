package soiregistry

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/koooyooo/soi-go/pkg/soi"

	"golang.org/x/xerrors"
)

type singleRegistry struct {
	soiCup *soi.SoiCup
}

func (sr singleRegistry) loaded() bool {
	return sr.soiCup != nil
}

// SoiCup returns SoiCup. when it is not loaded, or force is true, it load first.
func (sr *singleRegistry) SoiCup(force bool) (*soi.SoiCup, error) {
	if sr.loaded() && !force {
		return sr.soiCup, nil
	}
	return Load()
}

// Register registers a soi.go.
func (sr *singleRegistry) Register(s *soi.Soi) error {
	sc, err := sr.SoiCup(false)
	if err != nil {
		return err
	}
	newSois := soi.FilterByExcludeName(sc.SoisP(), s.Name)
	newSois = append(newSois, s)
	sc.SetP(newSois)
	err = Store(sc)
	if err != nil {
		return err
	}
	return nil
}

// findByFilter filters Sois by given function and return results
func (sr singleRegistry) findByFilter(f func(*soi.Soi) bool) ([]*soi.Soi, error) {
	sc, err := sr.SoiCup(false)
	if err != nil {
		return []*soi.Soi{}, err
	}
	var sois []*soi.Soi
	for _, s := range sc.SoisP() {
		if f(s) {
			sois = append(sois, s)
		}
	}
	return sois, nil
}

// FindByName returns Soi slice filtered by its name.
func (sr singleRegistry) FindByName(name string) ([]*soi.Soi, error) {
	return sr.findByFilter(func(s *soi.Soi) bool {
		return s.Name == name
	})
}

// FindByNamePart returns Soi slice filtered by its namePart.
func (sr singleRegistry) FindByNamePart(namePart string) ([]*soi.Soi, error) {
	return sr.findByFilter(func(s *soi.Soi) bool {
		return strings.Contains(s.Name, namePart)
	})
}

// FindByTag returns Soi slice filtered by its tags.
func (sr singleRegistry) FindByTag(tag string) ([]*soi.Soi, error) {
	return sr.findByFilter(func(s *soi.Soi) bool {
		for _, t := range s.Tags {
			if t == tag {
				return true
			}
		}
		return false
	})
}

//
func (sr *singleRegistry) Remove(name string) error {
	soiCup, err := sr.SoiCup(false)
	if err != nil {
		return err
	}
	removedSois := soi.FilterByExcludeName(soiCup.SoisP(), name)
	soiCup.SetP(removedSois)
	err = Store(soiCup)
	if err != nil {
		return err
	}
	return nil
}

func Load() (*soi.SoiCup, error) {
	soisFilePath, err := soi.SoisFilePath()
	if err != nil {
		return nil, err
	}
	s := soi.SoiCup{}
	b, err := ioutil.ReadFile(soisFilePath)
	if err != nil {
		return nil, xerrors.Errorf("error in reading [%s] %v", soi.SoisFilePath, err)
	}
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, xerrors.Errorf("error in unmarshalling content [%s] %v", string(b), err)
	}
	return &s, nil
}

func Store(s *soi.SoiCup) error {
	b, err := json.Marshal(s)
	if err != nil {
		return xerrors.Errorf("error in marshalling content %v %v", s, err)
	}
	var prettyBuff bytes.Buffer
	err = json.Indent(&prettyBuff, b, "", "  ")
	if err != nil {
		return xerrors.Errorf("filed in indent json %v", err)
	}
	soisFilePath, err := soi.SoisFilePath()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(soisFilePath, prettyBuff.Bytes(), 0666)
	if err != nil {
		return xerrors.Errorf("filed in writing file [%s] %v", soi.SoisFilePath, err)
	}
	return nil
}
