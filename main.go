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

func cookieTestHandleFunc(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	fmt.Fprintf(w, "Cookie set: %v", cookies)
}

func notFoundHandleFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Not Found</h1>")
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	us, err := models.NewUserService(psqlInfo)
	defer us.Close()
	if err != nil {
		panic(err)
	}
	us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/signup", usersC.SignUpView).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Auth).Methods("POST")
	r.HandleFunc("/cookietest", cookieTestHandleFunc).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(notFoundHandleFunc)
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
