package srv

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/koooyooo/soi-go/pkg/srv/repo"

	"github.com/koooyooo/soi-go/pkg/cli"
)

func Run() {
	m := http.NewServeMux()
	m.Handle("/", http.HandlerFunc(handler))
	m.Handle("/store", http.HandlerFunc(storeHandler))
	if err := http.ListenAndServe(":8080", m); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error @read-body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println(string(b)) // TODO
	var s cli.SoiWithPath
	if err = json.Unmarshal(b, &s); err != nil {
		log.Printf("error @unmarshal: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	ctx := context.Background()
	repo := repo.NewRepository()
	if err = repo.Store(ctx, &s); err != nil {
		log.Printf("error @store: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
