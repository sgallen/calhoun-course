package controllers

import (
	"fmt"
	"log"
	"net/http"

	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
)

// NewUser is used to create a new User controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers(us *models.UserService) *UserController {
	return &UserController{
		SignUpView: views.NewView(
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
	SignUpView  *views.View
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

// Create is used to process the signup form.
//
// POST /signup
func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("Route: %v", uc.SignUpView.Data.Route)
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

	err = uc.signIn(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (uc *UserController) Auth(w http.ResponseWriter, r *http.Request) {
	log.Printf("Route: %v", uc.LoginView.Data.Route)
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

	err = uc.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (uc *UserController) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		remember, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = remember

		err = uc.userService.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:  "remember_token",
		Value: user.Remember,
	}
	http.SetCookie(w, &cookie)

	return nil
}
