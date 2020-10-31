package srv

import (
	"fmt"
	"net/http"
)

func Run() {
	m := http.NewServeMux()
	m.Handle("/", http.HandlerFunc(handler))
	http.ListenAndServe(":8080", m)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}
