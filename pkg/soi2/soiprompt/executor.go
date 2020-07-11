package soiprompt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/koooyooo/soi-go/pkg/soi"
	"github.com/koooyooo/soi-go/pkg/soi2"
)

func Executor(in string) {
	fmt.Printf("EXEC: %s\n", in)
	cmd := strings.Split(in, " ")[0]
	subCmd := strings.TrimPrefix(in, cmd+" ")
	switch cmd {
	case "exit":
		os.Exit(0)
	case "open", "o", "list", "l":
		relPath := strings.ReplaceAll(subCmd, " ", "/")
		err := open(relPath)
		if err != nil {
			log.Fatal(err)
		}
	}
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
