package handler

import (
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	var templates = template.Must(template.ParseGlob("html/nested/*"))
	templates, err := templates.ParseFiles("html/" + tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, "base", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
