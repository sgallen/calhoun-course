package main

import (
	"html/template"
	"os"
)

type Dog struct {
	Name         string
	Breed        string
	FavoriteToys []string
}

type User struct {
	Name string
	Dog  Dog
	Age  int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	data := User{
		Name: "John Doe",
		Dog: Dog{
			Name:         "Fido",
			Breed:        "Lab",
			FavoriteToys: []string{"ball", "squishy", "blanket"},
		},
		Age: 35,
	}

	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
