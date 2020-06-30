/*
Package soicontrol offers control functions
*/
package soicontrol

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/koooyooo/soi-go/pkg/soi"

	"github.com/koooyooo/soi-go/pkg/soi/soiservice"
)

// Controller controls
type Controller struct {
	Service soiservice.SoiService
}

// Add is for adding new link to soi
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
		name = DefaultName(uri)
	}
	tags := TrimElements(strings.Split(*opTag, ","))
	soi, err := c.Service.Add(&soi.Soi{
		Name: name,
		Uri:  uri,
		Tags: tags,
	})
	if err != nil {
		log.Fatal("failed in adding", err)
	}
	fmt.Printf("added %v \n", soi)
}

// List is for listing up links
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
		sois = soi.FilterByTags(sois, tags)
	}
	showList(sois)
}

// Tags is for listing up all existing tags.
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
	sc := &soi.SoiCup{}
	sc.SetP(sois)
	tags := sc.TagSet(*namePart)
	sort.Strings(tags)
	showTags(tags)
}

// Open is for open a specified link
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

// Remove is for removing link
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
		// find exact name
		err := c.Service.Remove(soi.Name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("remove: %s \n", soi.Name)
	} else {
		// find no exact name and starts ambiguous search
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

// Tag is for adding tag for the link
func (c Controller) Tag() {
	flags := flag.NewFlagSet("tag", flag.PanicOnError)
	err := flags.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	name := flags.Arg(0)
	tagsStr := flags.Arg(1)
	tags := TrimElements(strings.Split(tagsStr, ","))
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

// ShowList outputs links which might be filtered
func showList(sois []*soi.Soi) {
	if len(sois) == 0 {
		fmt.Println("No Record Found")
	}
	for i, v := range sois {
		fmt.Printf(" - %02d:  %v\n", i+1, v)
	}
}

// ShowTags outputs tags
func showTags(tags []string) {
	if len(tags) == 0 {
		fmt.Println("No Record Found")
	}
	for i, v := range tags {
		fmt.Printf("- %02d:  %v\n", i+1, v)
	}
}
