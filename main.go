package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func prepareHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
}

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	prepareHeader(w)
	fmt.Fprint(w, "<h1>What's up?</h1>")
}

func contactHandleFunc(w http.ResponseWriter, r *http.Request) {
	prepareHeader(w)
	fmt.Fprint(
		w,
		`To get in touch, please send an email to
		<a href="mailto:foo@example.com">me</a>`,
	)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandleFunc)
	r.HandleFunc("/contact", contactHandleFunc)
	http.ListenAndServe(":3000", r)
}
