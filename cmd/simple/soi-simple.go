package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
)

// 標準入力の内容を読み込んで、SoiDataを生成して、Firefoxで開く
func main() {
	control()
}

type simpleSoiData struct {
	URI string `json:"uri"`
}

func control() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("fail in reading stdin: %s", err)
	}
	soi, err := loadSoi(b)
	if err != nil {
		log.Fatalf("fail in loading soi: %s", err)
	}
	openFirefox(soi)
}

func loadSoi(b []byte) (*simpleSoiData, error) {
	var sd simpleSoiData
	if err := json.Unmarshal(b, &sd); err != nil {
		return nil, err
	}
	return &sd, nil
}

// duplicated from pkg/cli/loader/loader.go
func openFirefox(s *simpleSoiData) error {
	return exec.Command("open", "-a", "Firefox", s.URI).Start()
}
