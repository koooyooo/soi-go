package execute

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/koooyooo/soi-go/pkg/common/hash"
	"github.com/koooyooo/soi-go/pkg/model"
	"github.com/koooyooo/soi-go/pkg/soiprompt/utils"
)

// add はsoiの追加を行います
func (e *Executor) add(in string) error {
	var tags utils.StringArray

	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	n := flags.String("n", "", "name of the uri")
	d := flags.String("d", "", "soiRoot to which model store")
	flags.Var(&tags, "t", "tag of the uri")
	if err := flags.Parse(strings.Split(in, " ")[1:]); err != nil {
		return err
	}

	dir, name, uri, parsedTags := parseArgs(flags.Args())
	if uri == "" {
		return fmt.Errorf("url is required")
	}
	if dir == "" {
		dir = time.Now().Format("2006-01")
	}
	if *n != "" {
		name = *n
	}
	if *d != "" {
		dir = *d
	}

	if name == "" {
		title, ok, err := utils.ParseTitleByURL(uri)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("no name(title) found in the url: %s\nplease specify the name with -n option\n", uri)
		}
		name = title
	}
	h, err := hash.Sha1(uri)
	if err != nil {
		return err
	}

	kTags, kvTags := separateTags(append([]string(tags), parsedTags...))

	s := &model.SoiData{
		Name:      name,
		Path:      dir,
		URI:       uri,
		Hash:      h,
		Tags:      kTags,
		KVTags:    kvTags,
		CreatedAt: time.Now(),
	}
	ctx := context.Background()
	if err := e.Service.Store(ctx, s); err != nil {
		return err
	}
	e.Cache.Clear()
	return nil
}

func parseArgs(args []string) (dir, name, uri string, tags []string) {
	var rests []string
	for _, a := range args {
		if strings.HasPrefix(a, "#") {
			tags = append(tags, strings.TrimPrefix(a, "#"))
			continue
		}
		if strings.HasPrefix(a, "http:") || strings.HasPrefix(a, "https:") {
			uri = a
			continue
		}
		rests = append(rests, a)
	}
	switch len(rests) {
	case 2:
		dir = rests[0]
		name = rests[1]
	case 1:
		dir = ""
		name = rests[0]
	default:
		dir = ""
		name = ""
	}
	return
}

// separateTags は入力されたタグをKタグとKVタグに仕分けます
func separateTags(tags []string) (kTags []string, kvTags []model.KVTag) {
	for _, tag := range tags {
		if strings.Contains(tag, "=") {
			kv := strings.SplitN(tag, "=", 2)
			kvTags = append(kvTags, model.KVTag{
				Key:   kv[0],
				Value: kv[1],
			})
		} else if tag != "" {
			kTags = append(kTags, tag)
		}
	}
	return
}
