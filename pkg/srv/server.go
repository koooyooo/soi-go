package srv

import (
	"fmt"
	"log"
	"net/http"
)

func Run() {
	m := http.NewServeMux()
	m.Handle("/", http.HandlerFunc(handler))
	if err := http.ListenAndServe(":8080", m); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}
