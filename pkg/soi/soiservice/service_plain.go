package soiservice

import (
	"github.com/koooyooo/soi-go/pkg/soi"
	"github.com/koooyooo/soi-go/pkg/soi/soiregistry"
)

type plainSoiService struct {
	Registry soiregistry.Registry
}

func (p plainSoiService) Add(s *soi.Soi) (*soi.Soi, error) {
	err := p.Registry.Register(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (p plainSoiService) Search(namePart string) ([]*soi.Soi, error) {
	return p.Registry.FindByNamePart(namePart)
}

func (p plainSoiService) SearchByTag(tag string) ([]*soi.Soi, error) {
	return p.Registry.FindByTag(tag)
}

func (p plainSoiService) Get(name string) (*soi.Soi, bool, error) {
	targetSois, err := p.Registry.FindByName(name)
	if err != nil {
		return nil, false, err
	}
	if 0 == len(targetSois) {
		return nil, false, nil
	}
	return targetSois[0], true, nil
}

func (p plainSoiService) Remove(name string) error {
	return p.Registry.Remove(name)
}

func (p plainSoiService) Tag(name string, tags []string) (*soi.Soi, bool, error) {
	soiBefore, ok, err := p.Get(name)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}

	soi, err := p.Add(&soi.Soi{
		Name: soiBefore.Name,
		Uri:  soiBefore.Uri,
		Tags: tags,
	})
	if err != nil {
		return nil, true, err
	}
	return soi, true, nil
}
