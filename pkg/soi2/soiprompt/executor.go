package soiprompt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/koooyooo/soi-go/pkg/soi"
	"github.com/koooyooo/soi-go/pkg/soi2"
)

func Executor(in string) {
	fmt.Printf("EXEC: %s\n", in)
	in = strings.Trim(in, " ")
	cmd := strings.Split(in, " ")[0]
	subCmd := strings.TrimPrefix(in, cmd+" ")
	switch cmd {
	case "exit":
		os.Exit(0)
	case "add":
		err := add(subCmd)
		if err != nil {
			log.Fatal(err)
		}
	case "open", "o", "list", "l":
		relPath := strings.ReplaceAll(subCmd, " ", "/")
		err := open(relPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func add(uri string) error {
	title, ok, err := parseTitleByURL(uri)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("no url found for: %s", uri)
	}
	s := soi2.SoiData{
		Name:    title,
		URI:     uri,
		Tags:    []string{},
		Created: fmt.Sprintf("%v", time.Now()),
	}
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	dir, err := soi.SoisDirPath()
	if err != nil {
		return err
	}
	baseDir := dir + "/new"
	if err = os.MkdirAll(baseDir, 0700); err != nil {
		return err
	}
	fileTitle := strings.ReplaceAll(title, " ", "_")
	return ioutil.WriteFile(baseDir+"/"+fileTitle+".json", b, 0600)
}

func open(relPath string) error {
	dir, err := soi.SoisDirPath()
	if err != nil {
		return err
	}
	fullPath := dir + "/" + relPath
	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	var soi soi2.SoiData
	err = json.Unmarshal(b, &soi)
	if err != nil {
		return err
	}
	err = exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", soi.URI).Start()
	if err != nil {
		return err
	}
	return nil
}
