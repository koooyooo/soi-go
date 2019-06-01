package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/koooyooo/soi-go/model"
	"github.com/koooyooo/soi-go/registory"
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
		err = Add(*name, uri, strings.Split(*tags, ","))
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

func Add(name, uri string, tags []string) error {
	fmt.Printf("soi add name=%s, url=%s \n", name, uri)
	sois, err := registory.Load()
	if err != nil {
		return err
	}
	if sois.Contains(name) {
		sois.Remove(name)
	}
	sois.Add(model.Soi{
		Name: name,
		Uri:  uri,
		Tags: tags,
	})

	fmt.Println(sois) // TODO
	err = registory.Store(*sois)
	if err != nil {
		return err
	}
	return nil
}
