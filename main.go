package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	name := flag.String("n", "", "name flag")
	tags := flag.String("t", "", "tag name")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Not enough argument")
	}
	cmd := args[0]
	switch cmd {
	case "add":
		err := Add(args, *name, strings.Split(*tags, ","))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Add(args []string, name string, tags []string) error {
	uri := args[1]
	fmt.Printf("soi add name=%s, url=%s \n", name, uri)
	sois, err := Load()
	if err != nil {
		return err
	}
	fmt.Println(sois)

	soi := Soi{
		Name: name,
		Uri:  uri,
		Tags: tags,
	}
	sois.Add(soi)
	fmt.Println(sois) // TODO
	err = Store(*sois)
	if err != nil {
		return err
	}
	return nil
}

type Sois struct {
	Sois []Soi `json:"sois"`
}

func (s *Sois) Add(soi Soi) {
	s.Sois = append(s.Sois, soi)
}

type Soi struct {
	Name string   `json:"name"`
	Uri  string   `json:"link"`
	Tags []string `json:"tags"`
}

func Load() (*Sois, error) {
	s := Sois{}
	b, err := ioutil.ReadFile("sois.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func Store(s Sois) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	var prettyBuff bytes.Buffer
	err = json.Indent(&prettyBuff, b, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("sois.json", prettyBuff.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}
