package opener

import (
	"github.com/koooyooo/soi-go/pkg/model"
	"os/exec"
)

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
