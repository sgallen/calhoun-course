package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var homeTemplate, contactTemplate *template.Template

func prepareHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	prepareHeader(w)
	if err := homeTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func contactHandleFunc(w http.ResponseWriter, r *http.Request) {
	prepareHeader(w)
	if err := contactTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func notFoundHandleFunc(w http.ResponseWriter, r *http.Request) {
	prepareHeader(w)
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Not Found</h1>")
}

func main() {
	var err error
	homeTemplate, err = template.ParseFiles("views/home.gohtml")
	if err != nil {
		panic(err)
	}
	contactTemplate, err = template.ParseFiles("views/contact.gohtml")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandleFunc)
	r.HandleFunc("/contact", contactHandleFunc)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandleFunc)
	http.ListenAndServe(":3000", r)
}
