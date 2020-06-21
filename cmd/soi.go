package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/koooyooo/soi-go/pkg/soi"

	"github.com/koooyooo/soi-go/pkg/soi/soiregistry"

	"github.com/koooyooo/soi-go/pkg/fileio"

	"github.com/koooyooo/soi-go/pkg/soi/soicontrol"

	"github.com/koooyooo/soi-go/pkg/soi/soiservice"
)

// main method of soi
func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Not enough argument")
	}
	registry, ok := soiregistry.NewRegistry(soiregistry.RegistryTypeLocal)
	if !ok {
		log.Fatal("no registry found")
	}
	service, ok := soiservice.NewSoiService(soiservice.ServiceTypePlain, registry)
	if !ok {
		log.Fatal("no service found")
	}
	controller := soicontrol.Controller{
		Service: service,
	}

	cmd := flag.Arg(0)
	switch cmd {
	case "i", "init":
		initSoi()
	case "a", "add":
		controller.Add()
	case "l", "list":
		controller.List()
	case "ts", "tags":
		controller.Tags()
	case "o", "open":
		controller.Open()
	case "r", "remove", "d", "delete":
		controller.Remove()
	case "t", "tag":
		controller.Tag()
	}
}

// create soi repository
func initSoi() {
	soisFilePath, err := soi.SoisFilePath()
	if err != nil {
		log.Fatal("failed in getting sois file path", err)
	}
	if fileio.FileExists(soisFilePath) {
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
