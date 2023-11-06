package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"soi-go/pkg/cli/constant"
	"soi-go/pkg/cli/soiprompt/utils"
	"soi-go/pkg/model"
	"strings"
)

func main() {
	fmt.Println("Hello Refresh Path")
	if err := control(); err != nil {
		log.Fatal(err)
	}
}

func control() error {
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
	filePath := strings.TrimSuffix(pathElms[len(pathElms)-1], ".json")

	fmt.Printf("%s: [%s]\n", dirPath, filePath)
	if soi.Path != dirPath {
		fmt.Printf("  [Fix] soi path is fixed with file path: [%s] -> [%s]\n", soi.Path, dirPath)
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
