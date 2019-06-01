package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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
		flags := flag.NewFlagSet("add", flag.PanicOnError)
		name := flags.String("n", "", "name flag")
		tags := flags.String("t", "", "tag name")
		err := flags.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
		uri := flags.Arg(0)
		if *name == "" {
			name = NameFromURI(uri)
		}
		err = service.Add(*name, uri, strings.Split(*tags, ","))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func NameFromURI(uri string) *string {
	idxURIServerStart := strings.Index(uri, "://") + len("://")
	fromStartUri := uri[idxURIServerStart:]
	idxFromStartServerEnd := strings.Index(fromStartUri, "/")
	if idxFromStartServerEnd == -1 {
		idxFromStartServerEnd = len(fromStartUri)
	}
	server := fromStartUri[:idxFromStartServerEnd]
	server = strings.TrimPrefix(server, "www.")
	return &server
}
