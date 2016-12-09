package handler

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseGlob("html/*"))

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
