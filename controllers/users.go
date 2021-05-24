package controllers

import (
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

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	log.Printf("Route: %v", u.View.Data.Route)
	if err := u.View.Render(w); err != nil {
		panic(err)
	}
}
