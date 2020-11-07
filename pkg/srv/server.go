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
	m.HandleFunc("/", handler)
	m.HandleFunc("/show", showHandler)
	m.HandleFunc("/store", storeHandler)
	if err := http.ListenAndServe(":8080", m); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	repo := repo.NewRepository()
	sb, err := repo.LoadAll(context.Background())
	if err != nil {
		log.Printf("error @show's loading : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(sb)
	if err != nil {
		log.Printf("error @show's marshaling: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		log.Printf("error @show's writting: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!!")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error @store's read-body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println(string(b)) // TODO
	var s cli.SoiWithPath
	if err = json.Unmarshal(b, &s); err != nil {
		log.Printf("error @store's unmarshal: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	repo := repo.NewRepository()
	if err = repo.Store(context.Background(), &s); err != nil {
		log.Printf("error @store's storing: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
