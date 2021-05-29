package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

func NewView(layout string, route string, files ...string) *View {
	addTemplateDir(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)

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

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := v.Render(w); err != nil {
		panic(err)
	}
}

func (v *View) Render(w http.ResponseWriter) error {
	return v.Template.ExecuteTemplate(w, v.Layout, v.Data)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

func addTemplateDir(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
