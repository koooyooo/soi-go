package main

import (
	"encoding/json"
	"fmt"
	"github.com/koooyooo/soi-go/pkg/cli/opener"
	"github.com/koooyooo/soi-go/pkg/model"
	"golang.org/x/term"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	control()
}

func control() {
	if term.IsTerminal(syscall.Stdin) {
		fmt.Println("This is a terminal")
		return
	}
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("fail in reading stdin: %s", err)
	}
	soi, err := loadSoi(b)
	if err != nil {
		log.Fatalf("fail in loading soi: %s", err)
	}
	opener.OpenFirefox(soi)
}

func loadSoi(b []byte) (*model.SoiData, error) {
	var sd model.SoiData
	if err := json.Unmarshal(b, &sd); err != nil {
		return nil, err
	}
	return &sd, nil
}

// duplicated from pkg/cli/loader/loader.go
func openFirefox(s *model.SoiData) error {
	return exec.Command("open", "-a", "Firefox", s.URI).Start()
}
