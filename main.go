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
	service := service.NewSoiService()

	cmd := flag.Arg(0)
	switch cmd {
	case "add":
		add(service)
	case "list":
		list(service)
	case "open":
		open(service)
	case "remove":
		remove(service)
	case "tag":
		tag(service)
	}
}

// add Link
func add(s service.SoiService) {
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
	soi, err := s.Add(name, uri, strings.Split(*tags, ","))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("added %v \n", soi)
}

func list(s service.SoiService) {
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
	sois, err := s.Search(*namePart)
	if err != nil {
		log.Fatal(err)
	}
	filteredSois := model.FilterByTags(sois, tags)
	showList(filteredSois)
}

func open(s service.SoiService) {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	namePart := flags.Arg(0)
	soi, ok, err := s.Get(namePart)
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

func remove(s service.SoiService) {
	flags := flag.NewFlagSet("remove", flag.PanicOnError)
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	namePart := flags.Arg(0)
	soi, ok, err := s.Get(namePart)
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		err := s.Remove(soi.Name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("remove: %s \n", soi.Name)
	} else {
		sois, err := s.Search(namePart)
		if err != nil {
			log.Fatal(err)
		}
		if len(sois) == 1 {
			err = s.Remove(sois[0].Name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("remove: %s \n", sois[0].Name)
		} else {
			showList(sois)
		}
	}
}

func tag(s service.SoiService) {
	flags := flag.NewFlagSet("tag", flag.PanicOnError)
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	name := flags.Arg(0)
	tagsStr := flags.Arg(1)
	tags := strings.Split(tagsStr, ",")
	soi, ok, err := s.Tag(name, tags)
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
