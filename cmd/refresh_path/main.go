package main

import (
	"encoding/json"
	"fmt"
	"github.com/koooyooo/soi-go/pkg/cli/constant"
	"github.com/koooyooo/soi-go/pkg/cli/soiprompt/utils"
	"github.com/koooyooo/soi-go/pkg/model"
	"log"
	"os"
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
		fmt.Printf("  [Warn] soi path is not matched with file path: [%s] [%s]\n", soi.Path, dirPath)
	}
	if soi.Name != filePath {
		fmt.Printf("  [Err] soi name is not matched with file name: [%s] [%s] [%s]\n", soi.Name, filePath, pathElms[0:1])
	}
	return nil
}
