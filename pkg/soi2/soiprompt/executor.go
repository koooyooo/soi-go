package soiprompt

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/koooyooo/soi-go/pkg/fileio"

	"github.com/koooyooo/soi-go/pkg/soi2"
)

func Executor(in string) {
	//fmt.Printf("EXEC: %s\n", in)
	in = strings.Trim(in, " ")
	cmd := strings.Split(in, " ")[0]
	switch cmd {
	case "quit", "q", "exit":
		err := quit(in)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "add", "a":
		err := add(in)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "mv":
		err := mv(in)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "rm":
		err := rm(in)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "open", "o", "list", "l", "dig", "d":
		err := open(in)
		if err != nil {
			fmt.Println(err)
			return
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
			return fmt.Errorf("no name(title) found in the url: %s\nplease specify the name with -n option\n", uri)
		}
		name = title
	}

	s := soi2.SoiData{
		Name:    name,
		URI:     uri,
		Tags:    []string{},
		Created: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	soiRoot, err := fileio.SoisDirPath()
	if err != nil {
		return err
	}
	baseDir := filepath.Join(soiRoot, *d)
	if err = os.MkdirAll(baseDir, 0700); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(baseDir, toStorableName(name)), b, 0600)
}

func mv(in string) error {
	baseDir, err := fileio.SoisDirPath()
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("mv", flag.PanicOnError)
	flags.Parse(strings.Split(in, " ")[1:])
	from := filepath.Join(baseDir, flags.Arg(0))
	to := filepath.Join(baseDir, flags.Arg(1))

	// 移動先のディレクトリが存在しない場合は作成
	toDir := to[0:strings.LastIndex(to, "/")]
	if !fileio.FileExists(toDir) {
		err = os.MkdirAll(toDir, 0700)
		if err != nil {
			return err
		}
	}

	// 末尾JSONの付与
	toIsDir, err := fileio.IsDir(to)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
		default:
			return err
		}
	}
	if !toIsDir && !strings.HasSuffix(to, ".json") {
		to = to + ".json"
	}
	return exec.Command("mv", from, to).Start()
}

func rm(in string) error {
	baseDir, err := fileio.SoisDirPath()
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("rm", flag.PanicOnError)
	flags.Parse(strings.Split(in, " ")[1:])

	target := filepath.Join(baseDir, flags.Arg(0))
	if !fileio.FileExists(target) {
		fmt.Println("No file or dir found.")
		return nil
	}
	return exec.Command("rm", "-rf", target).Start()
}

func open(in string) error {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	//chrome := flags.Bool("c", false, "use chrome")
	firefox := flags.Bool("f", false, "use firefox")
	safari := flags.Bool("s", false, "use safari")
	flags.Parse(strings.Split(in, " ")[1:])

	soisDir, err := fileio.SoisDirPath()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(soisDir, filepath.Join(flags.Args()...))
	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	var soi soi2.SoiData
	err = json.Unmarshal(b, &soi)
	if err != nil {
		return err
	}
	if *firefox {
		return exec.Command("open", "-a", "Firefox", soi.URI).Start()
	}
	if *safari {
		return exec.Command("open", "-a", "Safari", soi.URI).Start()
	}
	return exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", soi.URI).Start()
}

func quit(in string) error {
	fmt.Println("bye!!")
	os.Exit(0)
	return nil
}
