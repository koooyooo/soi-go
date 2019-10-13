package registory

import (
	"github.com/koooyooo/soi-go/model"
)

type Registory interface {
	Load() (*model.SoiCup, error)
	Store(s model.SoiCup) error
}

func NewRegistory() Registory {
	return localRegistroy{}
}
