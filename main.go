package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
