package main

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/controllers"
	"lenslocked.com/views"

	"github.com/gorilla/mux"
)

func init() {
	log.SetPrefix("MAIN: ")
}

var (
	homeView    *views.View
	contactView *views.View
)

func prepareHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
	prepareHeader(w)
	log.Printf("Route: %v", homeView.Data.Route)
	must(homeView.Render(w))
}

func contactHandleFunc(w http.ResponseWriter, r *http.Request) {
	prepareHeader(w)
	log.Printf("Route: %v", contactView.Data.Route)
	must(contactView.Render(w))
}

func notFoundHandleFunc(w http.ResponseWriter, r *http.Request) {
	prepareHeader(w)
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Not Found</h1>")
}

func main() {
	homeView = views.NewView("bootstrap", "home", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "contact", "views/contact.gohtml")

	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandleFunc)
	r.HandleFunc("/contact", contactHandleFunc)
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.NotFoundHandler = http.HandlerFunc(notFoundHandleFunc)
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
