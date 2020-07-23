package soiprompt

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/koooyooo/soi-go/pkg/fileio"

	"github.com/koooyooo/soi-go/pkg/soi"
	"github.com/koooyooo/soi-go/pkg/soi2"
)

func Executor(in string) {
	//fmt.Printf("EXEC: %s\n", in)
	in = strings.Trim(in, " ")
	cmd := strings.Split(in, " ")[0]
	subCmd := strings.TrimPrefix(in, cmd+" ")
	switch cmd {
	case "quit", "q", "exit":
		fmt.Println("quit")
		os.Exit(0)
	case "add", "a":
		err := add(in)
		if err != nil {
			log.Fatal(err)
		}
	case "mv":
		err := mv(in)
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

func add(in string) error {
	flags := flag.NewFlagSet("add", flag.PanicOnError)
	n := flags.String("n", "", "name of the uri")
	d := flags.String("d", "new", "soiRoot to which soi store")
	err := flags.Parse(strings.Split(in, " ")[1:])
	if err != nil {
		return err
	}

	uri := flags.Arg(0)

	name := *n
	if name == "" {
		title, ok, err := parseTitleByURL(uri)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("no url found for: %s", uri)
		}
		name = title
	}

	s := soi2.SoiData{
		Name:    name,
		URI:     uri,
		Tags:    []string{},
		Created: fmt.Sprintf("%v", time.Now()),
	}
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	soiRoot, err := soi.SoisDirPath()
	if err != nil {
		return err
	}
	baseDir := filepath.Join(soiRoot, *d)
	if err = os.MkdirAll(baseDir, 0700); err != nil {
		return err
	}
	fileName := strings.ReplaceAll(name, " ", "_")
	return ioutil.WriteFile(baseDir+"/"+fileName+".json", b, 0600)
}

func mv(in string) error {
	baseDir, err := soi.SoisDirPath()
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("mv", flag.PanicOnError)
	flags.Parse(strings.Split(in, " ")[1:])
	from := baseDir + "/" + strings.TrimPrefix(flags.Arg(0), "/")
	to := baseDir + "/" + strings.TrimPrefix(flags.Arg(1), "/")

	toDir := to[0:strings.LastIndex(to, "/")]
	if !fileio.FileExists(toDir) {
		err = os.Mkdir(toDir, 0700)
		if err != nil {
			return err
		}
	}
	return exec.Command("mv", from, to).Start()
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
	return exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", soi.URI).Start()
}
