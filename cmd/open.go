/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "open input files",
	Long:  `Opens the specified URI in Firefox browser. ts JSON input from stdin with a "uri" field.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := open(); err != nil {
			log.Fatalf("failed to open: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}

type simpleSoiData struct {
	URI string `json:"uri"`
}

// open reads JSON from stdin and opens the URI in Firefox
func open() error {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("fail in reading stdin: %w", err)
	}
	soi, err := loadSoi(b)
	if err != nil {
		return fmt.Errorf("fail in loading soi: %w", err)
	}
	if err := openFirefox(soi); err != nil {
		return fmt.Errorf("fail in opening browser: %w", err)
	}
	return nil
}

// loadSoi loads the simpleSoiData from the input bytes
func loadSoi(b []byte) (*simpleSoiData, error) {
	var sd simpleSoiData
	if err := json.Unmarshal(b, &sd); err != nil {
		return nil, err
	}
	if !strings.HasPrefix(sd.URI, "http") {
		return nil, fmt.Errorf("invalid uri: %s", sd.URI)
	}
	return &sd, nil
}

// duplicated from pkg/cli/loader/loader.go
func openFirefox(s *simpleSoiData) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", "-a", "Firefox", s.URI).Start()
	case "linux":
		return exec.Command("firefox", s.URI).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", "firefox", s.URI).Start()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
