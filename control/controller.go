package control

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/koooyooo/soi-go/common"
	"github.com/koooyooo/soi-go/model"
	"github.com/koooyooo/soi-go/service"
)

// Controller controls
type Controller struct {
	Service service.SoiService
}

// add is for adding new link to soi
func (c Controller) Add() {
	flags := flag.NewFlagSet("add", flag.PanicOnError)
	opTag := flags.String("t", "", "tags")
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal("failed in parsing args", err)
	}
	name := flags.Arg(0)
	uri := flags.Arg(1)
	if name == "-" {
		name = common.DefaultName(uri)
	}
	tags := common.TrimElements(strings.Split(*opTag, ","))
	soi, err := c.Service.Add(name, uri, tags)
	if err != nil {
		log.Fatal("failed in adding", err)
	}
	fmt.Printf("added %v \n", soi)
}

// list is for listing up links
func (c Controller) List() {
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
	sois, err := c.Service.Search(*namePart)
	if err != nil {
		log.Fatal("failed in searching", err)
	}
	if len(tags) > 0 {
		sois = model.FilterByTags(sois, tags)
	}
	showList(sois)
}

// tags
func (c Controller) Tags() {
	flags := flag.NewFlagSet("tags", flag.PanicOnError)
	namePart := flags.String("n", "", "name")

	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal("failed in parsing", err)
	}
	sois, err := c.Service.Search("")
	if err != nil {
		log.Fatal("failed in searching", err)
	}
	tags := (&model.SoiCup{Sois: sois}).TagSet(*namePart)
	sort.Strings(tags)
	showTags(tags)
}

// open is for open a specified link
func (c Controller) Open() {
	flags := flag.NewFlagSet("open", flag.PanicOnError)
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	namePart := flags.Arg(0)
	soi, ok, err := c.Service.Get(namePart)
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
func (c Controller) Remove() {
	flags := flag.NewFlagSet("remove", flag.PanicOnError)
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	namePart := flags.Arg(0)
	soi, ok, err := c.Service.Get(namePart)
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		err := c.Service.Remove(soi.Name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("remove: %s \n", soi.Name)
	} else {
		sois, err := c.Service.Search(namePart)
		if err != nil {
			log.Fatal(err)
		}
		if len(sois) == 1 {
			err = c.Service.Remove(sois[0].Name)
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
func (c Controller) Tag() {
	flags := flag.NewFlagSet("tag", flag.PanicOnError)
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	name := flags.Arg(0)
	tagsStr := flags.Arg(1)
	tags := common.TrimElements(strings.Split(tagsStr, ","))
	soi, ok, err := c.Service.Tag(name, tags)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		fmt.Println("no target found")
		return
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
