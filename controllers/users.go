package controllers

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/views"
)

// NewUser is used to create a new User controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers(us *models.UserService) *UserController {
	return &UserController{
		View: views.NewView(
			"bootstrap",
			"signup",
			"users/new",
		),
		LoginView: views.NewView(
			"bootstrap",
			"login",
			"users/auth",
		),
		userService: us,
	}
}

type UserController struct {
	View        *views.View
	LoginView   *views.View
	userService *models.UserService
}

type SignUpForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
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
func (uc *UserController) New(w http.ResponseWriter, r *http.Request) {
	log.Printf("Route: %v", uc.View.Data.Route)
	if err := uc.View.Render(w); err != nil {
		panic(err)
	}
}

// Create is used to process the signup form.
//
// POST /signup
func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var form SignUpForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	err := uc.userService.Create(&user)
	switch err {
	case models.ErrInvalidPassword:
		fmt.Fprintln(w, "Invalid password")
	case nil:
		fmt.Fprintf(w, "SignUpForm struct: %v", form)
		fmt.Fprintf(w, "Created user: %v", user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Route: %v", uc.View.Data.Route)
	if err := uc.LoginView.Render(w); err != nil {
		panic(err)
	}
}

func (uc *UserController) Auth(w http.ResponseWriter, r *http.Request) {
	log.Printf("Route: %v", uc.View.Data.Route)
	var form LoginForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := uc.userService.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address")
			return
		case models.ErrIncorrectPassword:
			fmt.Fprintln(w, "Incorrect password")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	createCookie(w, user)
	fmt.Fprintln(w, user)
}

func createCookie(w http.ResponseWriter, user *models.User) {
	cookie := http.Cookie{
		Name:  "email",
		Value: user.Email,
	}
	http.SetCookie(w, &cookie)
}
