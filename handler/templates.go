package handler

import (
	"html/template"
	"net/http"

	"strings"

	"github.com/martinhoefling/go-dmarc-report/bindata"
)

var parsedTemplates map[string]*template.Template

func getTemplate(templateName string) *template.Template {
	if parsedTemplates == nil {
		parsedTemplates = make(map[string]*template.Template)
	}

	tmpl, ok := parsedTemplates[templateName]
	if ok {
		return tmpl
	}

	first := true
	var newTmpl *template.Template
	for _, name := range bindata.AssetNames() {
		if strings.HasPrefix(name, "html/nested/") {
			var newTemplate *template.Template
			if first {
				newTemplate = template.New(name)
				first = false
			} else {
				newTemplate = newTmpl.New(name)
			}
			templateString := string(bindata.MustAsset(name))
			newTmpl = template.Must(newTemplate.Parse(templateString))
		}
	}
	templateString := string(bindata.MustAsset("html/" + templateName + ".html"))
	newTemplate := newTmpl.New(templateName)
	newTmpl = template.Must(newTemplate.Parse(templateString))

	parsedTemplates[templateName] = newTmpl
	return parsedTemplates[templateName]
}

func renderTemplate(w http.ResponseWriter, templateName string, p interface{}) {
	tmpl := getTemplate(templateName)
	err := tmpl.ExecuteTemplate(w, "base", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
