package macos

import (
	"errors"
	"os/exec"
	"github.com/koooyooo/soi-go/pkg/cli/opener"
	"github.com/koooyooo/soi-go/pkg/model"
)

func NewOpener() opener.Opener {
	return &macOSOpener{}
}

type macOSOpener struct{}

func (o macOSOpener) OpenChrome(s *model.SoiData, _ bool) error {
	return exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", s.URI).Start()
}

func (o macOSOpener) OpenFirefox(s *model.SoiData, private bool) error {
	if private {
		return exec.Command("/Applications/Firefox.app/Contents/MacOS/firefox-bin", "-private-window", s.URI).Start()
	}
	return exec.Command("open", "-a", "Firefox", s.URI).Start()
}

func (o macOSOpener) OpenSafari(s *model.SoiData, _ bool) error {
	return exec.Command("open", "-a", "Safari", s.URI).Start()
}

func (o macOSOpener) OpenEdge(s *model.SoiData, _ bool) error {
	return errors.New("not implemented")
}

func OpenChrome(s *model.SoiData, _ bool) error {
	return exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", s.URI).Start()
}

func OpenFirefox(s *model.SoiData, private bool) error {
	if private {
		return exec.Command("/Applications/Firefox.app/Contents/MacOS/firefox-bin", "-private-window", s.URI).Start()
	}
	return exec.Command("open", "-a", "Firefox", s.URI).Start()
}

func OpenSafari(s *model.SoiData, _ bool) error {
	return exec.Command("open", "-a", "Safari", s.URI).Start()
}
