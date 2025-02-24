package opener

import (
	"errors"
	"os/exec"

	"github.com/koooyooo/soi-go/pkg/model"
)

func NewWindowsOpener() Opener {
	return &winOpener{}
}

type winOpener struct{}

func (o winOpener) OpenChrome(s *model.SoiData, private bool) error {
	if private {
		return exec.Command("cmd", "/c", "start", "chrome", "--incognito", s.URI).Start()
	}
	return exec.Command("cmd", "/c", "start", "chrome", s.URI).Start()
}

func (o winOpener) OpenFirefox(s *model.SoiData, private bool) error {
	if private {
		return exec.Command("cmd", "/c", "start", "firefox", "-private-window", s.URI).Start()
	}
	return exec.Command("cmd", "/c", "start", "firefox", s.URI).Start()
}

func (o winOpener) OpenSafari(s *model.SoiData, _ bool) error {
	return errors.New("safari is not available on windows")
}

func (o winOpener) OpenEdge(s *model.SoiData, private bool) error {
	if private {
		return exec.Command("cmd", "/c", "start", "msedge", "--inprivate", s.URI).Start()
	}
	return exec.Command("cmd", "/c", "start", "msedge", s.URI).Start()
}
