package controllers

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/views"
)

// NewUser is used to create a new User controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers() *Users {
	return &Users{
		View: views.NewView(
			"bootstrap",
			"signup",
			"views/users/new.gohtml",
		),
	}
}

type Users struct {
	View *views.View
}

// New is used to render the form where a user can create
// a new user account.
//
// GET /signup
//
// TODO:
// I don't like the design Calhoun is using here. Would prefer
// the signup page to be a template that's focused on rendering
// an HTML page and then have standard REST endpoints:
// GET /users - fetch all users
// POST /users - create a user
// GET /users/<id> fetch a user
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	log.Printf("Route: %v", u.View.Data.Route)
	if err := u.View.Render(w); err != nil {
		panic(err)
	}
}

// Create is used to process the signup form.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a temp response.")
}
