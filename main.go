package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func prepareHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
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

func notFoundHandleFunc(w http.ResponseWriter, r *http.Request) {
	prepareHeader(w)
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Not Found</h1>")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandleFunc)
	r.HandleFunc("/contact", contactHandleFunc)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandleFunc)
	http.ListenAndServe(":3000", r)
}
