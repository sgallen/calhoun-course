package controllers

import (
	"lenslocked.com/views"
)

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "home", "static/home"),
		Contact: views.NewView("bootstrap", "contact", "static/contact"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
}
