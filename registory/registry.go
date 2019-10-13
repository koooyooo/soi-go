package registory

import (
	"github.com/koooyooo/soi-go/model"
)

type Registry interface {
	Load() (*model.SoiCup, error)
	Store(s model.SoiCup) error
}

func NewRegistry() Registry {
	return localRegistry{}
}
