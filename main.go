package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	switch page := r.URL.Path; page {
	case "/":
		fmt.Fprint(w, "<h1>What's up?</h1>")
	case "/contact":
		fmt.Fprint(
			w,
			`To get in touch, please send an email to
			<a href="mailto:foo@example.com">me</a>`,
		)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
