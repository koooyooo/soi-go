package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/koooyooo/soi-go/model"

	util "github.com/koooyooo/soi-go/comons/uri"

	"github.com/koooyooo/soi-go/service"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Not enough argument")
	}
	cmd := flag.Arg(0)
	switch {
	case cmd == "a" || cmd == "add":
		add()
	case cmd == "l" || cmd == "list":
		list()
	case cmd == "o" || cmd == "open":
		open()
	case cmd == "r" || cmd == "remove":
		remove()
	case cmd == "t" || cmd == "tag":
		tag()
	}
}

func add() {
	flags := flag.NewFlagSet("add", flag.PanicOnError)
	tags := flags.String("t", "", "tags")
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	name := flags.Arg(0)
	uri := flags.Arg(1)
	if name == "-" {
		name = util.DefaultName(uri)
	}
	soi, err := service.Add(name, uri, strings.Split(*tags, ","))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("added %v \n", soi)
}

func list() {
	flags := flag.NewFlagSet("list", flag.PanicOnError)
	namePart := flags.String("n", "", "name")
	tagStr := flags.String("t", "", "tags")

	err := flags.Parse(os.Args[2:])
	var tags []string
	if *tagStr != "" {
		tags = strings.Split(*tagStr, ",")
	}
	if err != nil {
		log.Fatal(err)
	}
	sois, err := service.Search(*namePart)
	if err != nil {
		log.Fatal(err)
	}
	filteredSois := model.FilterByTags(sois, tags)
	showList(filteredSois)
}

func open() {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	namePart := flags.Arg(0)
	soi, ok, err := service.Get(namePart)
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

func remove() {
	flags := flag.NewFlagSet("remove", flag.PanicOnError)
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	namePart := flags.Arg(0)
	soi, ok, err := service.Get(namePart)
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		err := service.Remove(soi.Name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("remove: %s \n", soi.Name)
	} else {
		sois, err := service.Search(namePart)
		if err != nil {
			log.Fatal(err)
		}
		if len(sois) == 1 {
			err = service.Remove(sois[0].Name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("remove: %s \n", sois[0].Name)
		} else {
			showList(sois)
		}
	}
}

func tag() {
	flags := flag.NewFlagSet("tag", flag.PanicOnError)
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	name := flags.Arg(0)
	tagsStr := flags.Arg(1)
	tags := strings.Split(tagsStr, ",")
	soi, ok, err := service.Tag(name, tags)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		fmt.Println("no target found")
	}
	fmt.Printf("%v \n", soi)
}

func showList(sois []model.Soi) {
	if len(sois) == 0 {
		fmt.Println("No Record Found")
	}
	for i, v := range sois {
		fmt.Printf(" - %02d:  %v", i+1, v)
	}
}
