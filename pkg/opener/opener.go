package opener

import "github.com/koooyooo/soi-go/pkg/model"

type Opener interface {
	OpenChrome(*model.SoiData, bool) error
	OpenFirefox(*model.SoiData, bool) error
	OpenSafari(*model.SoiData, bool) error
	OpenEdge(*model.SoiData, bool) error
}
