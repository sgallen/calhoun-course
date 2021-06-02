package main

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/controllers"
	"lenslocked.com/models"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "foo"
	password = "foo"
	dbname   = "usegolang_dev"
)

func init() {
	log.SetPrefix("MAIN: ")
}

func notFoundHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Not Found</h1>")
}

func main() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	// us.DestructiveReset()
	user, err := us.ById(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

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
