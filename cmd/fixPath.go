/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"soi-go/pkg/cli/constant"
	"soi-go/pkg/cli/soiprompt/utils"
	"soi-go/pkg/model"
	"strings"

	"github.com/spf13/cobra"
)

// fixPathCmd represents the fixPath command
var fixPathCmd = &cobra.Command{
	Use:   "fixPath",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := controlFixPath(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(fixPathCmd)
}

func controlFixPath() error {
	soisDir, err := constant.SoisDir()
	if err != nil {
		return err
	}
	paths, err := utils.ListFilesRecursively(soisDir)
	if err != nil {
		return err
	}
	for _, p := range paths {
		if err := controlFile(soisDir, p); err != nil {
			fmt.Fprintf(os.Stderr, "failed in control file: %s", p)
		}
	}
	return nil
}

func controlFile(soisDir, p string) error {
	if strings.HasSuffix(p, "config.json") {
		return nil
	}
	if !strings.HasSuffix(p, ".json") {
		return nil
	}
	d, err := os.ReadFile(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed in reading file: %s", p)
	}
	var soi model.SoiData
	if err := json.Unmarshal(d, &soi); err != nil {
		fmt.Fprintf(os.Stderr, "failed in unmarshaling soi: %s", p)
	}
	pathFromSoiDir := strings.TrimPrefix(p, soisDir)
	pathFromSoiDir = strings.Trim(pathFromSoiDir, "/")

	pathElms := strings.Split(pathFromSoiDir, "/")
	dirPath := strings.Join(pathElms[1:len(pathElms)-1], "/") // skipping bucket elm from head & file one from tail
	fileName := pathElms[len(pathElms)-1]
	soiName := strings.TrimSuffix(fileName, ".json")

	if soi.Name != soiName {
		soi.Name = soiName
	}
	if soi.Path != dirPath {
		fmt.Printf("  [Fix] soi path for [%s] is fixed with file path:  [%s] -> [%s]\n", soiName, soi.Path, dirPath)
		soi.Path = dirPath
	}
	d, err = json.Marshal(soi)
	if err != nil {
		return err
	}
	if err := os.WriteFile(p, d, 644); err != nil {
		return err
	}
	return nil
}
