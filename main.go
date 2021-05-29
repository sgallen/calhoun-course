package main

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/controllers"

	"github.com/gorilla/mux"
)

func init() {
	log.SetPrefix("MAIN: ")
}

func notFoundHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Not Found</h1>")
}

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
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
