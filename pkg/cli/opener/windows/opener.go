package windows

import (
	"errors"
	"github.com/koooyooo/soi-go/pkg/cli/opener"
	"github.com/koooyooo/soi-go/pkg/model"
	"os/exec"
)

func NewOpener() opener.Opener {
	return &winOpener{}
}

type winOpener struct{}

func (o winOpener) OpenChrome(s *model.SoiData, _ bool) error {
	return exec.Command("start", "firefox", s.URI).Start()
}

func (o winOpener) OpenFirefox(s *model.SoiData, private bool) error {
	if private {
		return exec.Command("start", "firefox", "-private-window", s.URI).Start()
	}
	return exec.Command("start", "firefox", s.URI).Start()
}

func (o winOpener) OpenSafari(s *model.SoiData, _ bool) error {
	return errors.New("not implemented")
}

func (o winOpener) OpenEdge(s *model.SoiData, _ bool) error {
	return exec.Command("start", "msedge", s.URI).Start()
}
