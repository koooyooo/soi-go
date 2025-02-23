package wsl

import (
	"os/exec"

	"github.com/koooyooo/soi-go/pkg/cli/opener"
	"github.com/koooyooo/soi-go/pkg/model"
)

type wslOpener struct{}

func NewOpener() opener.Opener {
	return &wslOpener{}
}

func (o wslOpener) OpenChrome(s *model.SoiData, _ bool) error {
	return exec.Command("cmd.exe", "/C", "start", "chrome", s.URI).Start()
}

func (o wslOpener) OpenFirefox(s *model.SoiData, private bool) error {
	if private {
		return exec.Command("cmd.exe", "/C", "start", "firefox", "-private-window", s.URI).Start()
	}
	return exec.Command("cmd.exe", "/C", "start", "firefox", s.URI).Start()
}

func (o wslOpener) OpenSafari(s *model.SoiData, _ bool) error {
	// Safari is not available on Windows/WSL
	return exec.Command("cmd.exe", "/C", "start", s.URI).Start()
}

func (o wslOpener) OpenEdge(s *model.SoiData, _ bool) error {
	return exec.Command("cmd.exe", "/C", "start", "msedge", s.URI).Start()
}
