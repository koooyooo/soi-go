/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "open input files",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		open()
	},
}

func init() {
	rootCmd.AddCommand(openCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type simpleSoiData struct {
	URI string `json:"uri"`
}

func open() {
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
