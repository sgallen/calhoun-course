package controllers

import (
	"lenslocked.com/views"
)

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "home", "views/static/home.gohtml"),
		Contact: views.NewView("bootstrap", "contact", "views/static/contact.gohtml"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
}
