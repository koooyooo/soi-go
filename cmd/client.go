package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/koooyooo/soi-go/comons"

	"github.com/koooyooo/soi-go/model"

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
	case "init":
		initSoi()
	case "a", "add":
		add(service)
	case "l", "list":
		list(service)
	case "ts", "tags":
		tags(service)
	case "o", "open":
		open(service)
	case "r", "remove", "d", "delete":
		remove(service)
	case "t", "tag":
		tag(service)
	}
}

func initSoi() {
	soisFilePath, err := commons.SoisFilePath()
	if err != nil {
		log.Fatal("failed in getting sois file path", err)
	}
	if commons.FileExists(soisFilePath) {
		r, err := ioutil.ReadFile(soisFilePath)
		if err != nil {
			log.Fatal("failed in reading sois file", err)
		}
		ioutil.WriteFile(soisFilePath+".bk", r, 0600)
		if err := os.Remove(soisFilePath); err != nil {
			log.Fatal("failed in removing old sois file", err)
		}
	}
	ioutil.WriteFile(soisFilePath, []byte(`{"sois": []}`), 0600)
}

// add is for adding new link to soi
func add(s service.SoiService) {
	flags := flag.NewFlagSet("add", flag.PanicOnError)
	opTag := flags.String("t", "", "tags")
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal("failed in parsing args", err)
	}
	name := flags.Arg(0)
	uri := flags.Arg(1)
	if name == "-" {
		name = commons.DefaultName(uri)
	}
	tags := commons.TrimElements(strings.Split(*opTag, ","))
	soi, err := s.Add(name, uri, tags)
	if err != nil {
		log.Fatal("failed in adding", err)
	}
	fmt.Printf("added %v \n", soi)
}

// list is for listing up links
func list(s service.SoiService) {
	flags := flag.NewFlagSet("list", flag.PanicOnError)
	namePart := flags.String("n", "", "name")
	tagStr := flags.String("t", "", "tags")

	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatalf("failed in parsing args: %+v", err)
	}
	var tags []string
	if *tagStr != "" {
		tags = strings.Split(*tagStr, ",")
	}
	sois, err := s.Search(*namePart)
	if err != nil {
		log.Fatal("failed in searching", err)
	}
	if len(tags) > 0 {
		sois = model.FilterByTags(sois, tags)
	}
	showList(sois)
}

// tags
func tags(s service.SoiService) {
	flags := flag.NewFlagSet("tags", flag.PanicOnError)
	namePart := flags.String("n", "", "name")

	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal("failed in parsing", err)
	}
	sois, err := s.Search("")
	if err != nil {
		log.Fatal("failed in searching", err)
	}
	tags := (&model.SoiCup{Sois: sois}).TagSet(*namePart)
	sort.Strings(tags)
	showTags(tags)
}

// open is for open a specified link
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

// remove is for removing link
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

// tag is for adding tag for the link
func tag(s service.SoiService) {
	flags := flag.NewFlagSet("tag", flag.PanicOnError)
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	name := flags.Arg(0)
	tagsStr := flags.Arg(1)
	tags := commons.TrimElements(strings.Split(tagsStr, ","))
	soi, ok, err := s.Tag(name, tags)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		fmt.Println("no target found")
	}
	fmt.Printf("%v \n", soi)
}

// showList outputs links which might be filtered
func showList(sois []model.Soi) {
	if len(sois) == 0 {
		fmt.Println("No Record Found")
	}
	for i, v := range sois {
		fmt.Printf(" - %02d:  %v\n", i+1, v)
	}
}

func showTags(tags []string) {
	if len(tags) == 0 {
		fmt.Println("No Record Found")
	}
	for i, v := range tags {
		fmt.Printf("- %02d:  %v\n", i+1, v)
	}
}
