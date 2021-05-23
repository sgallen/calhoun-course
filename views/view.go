package views

import "html/template"

func NewView(layout string, route string, files ...string) *View {
	files = append(
		files,
		"views/layouts/bootstrap.gohtml",
		"views/layouts/navbar.gohtml",
		"views/layouts/footer.gohtml",
	)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
		Data:     Data{Route: route},
	}
}

type Data struct {
	Route string
}

type View struct {
	Template *template.Template
	Layout   string
	Data     Data
}
