package linux

import (
	"os/exec"
	"soi-go/pkg/cli/opener"
	"soi-go/pkg/model"
)

type linuxOpener struct{}

func NewLinuxOpener() opener.Opener {
	return &linuxOpener{}
}

func (o linuxOpener) OpenChrome(s *model.SoiData, _ bool) error {
	return exec.Command("google-chrome", s.URI).Start()
}

func (o linuxOpener) OpenFirefox(s *model.SoiData, private bool) error {
	if private {
		return exec.Command("firefox", "-private-window", s.URI).Start()
	}
	return exec.Command("firefox", s.URI).Start()
}

func (o linuxOpener) OpenSafari(s *model.SoiData, _ bool) error {
	return exec.Command("safari", s.URI).Start()
}

func (o linuxOpener) OpenEdge(s *model.SoiData, _ bool) error {
	return exec.Command("msedge", s.URI).Start()
}
