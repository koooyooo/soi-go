package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	util "github.com/koooyooo/soi-go/comons/uri"

	"github.com/koooyooo/soi-go/service"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Not enough argument")
	}
	cmd := flag.Arg(0)
	switch cmd {
	case "add":
		add()
	case "list":
		list()
	case "open":
		open()
	}
}

func add() {
	flags := flag.NewFlagSet("add", flag.PanicOnError)
	name := flags.String("n", "", "name flag")
	tags := flags.String("t", "", "tag name")
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	uri := flags.Arg(0)
	if *name == "" {
		name = util.DefaultName(uri)
	}
	err = service.Add(*name, uri, strings.Split(*tags, ","))
	if err != nil {
		log.Fatal(err)
	}
}

func list() {
	flags := flag.NewFlagSet("list", flag.PanicOnError)
	namePart := flags.String("n", "", "name")
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	sois, err := service.Search(*namePart)
	if err != nil {
		log.Fatal(err)
	}
	for i, v := range sois {
		fmt.Printf(" - %02d:  %-18s %s\n", i+1, v.Name, v.Uri)
	}
}

func open() {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	namePart := flags.String("n", "", "name")
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	soi, ok, err := service.Get(*namePart)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		fmt.Println("no data found")
		return
	}
	err = exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", soi.Uri).Start()
	if err != nil {
		log.Fatal(err)
	}
}
