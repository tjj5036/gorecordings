package util

import (
	"html/template"
	"log"
	"net/http"
)

// RenderTemplate renders a template given a template file to render and data
// in the form of an interface
func RenderTemplate(w http.ResponseWriter, template_file string, data interface{}) {
	var tpl *template.Template
	tpl = template.Must(template.ParseGlob("templates/*html"))
	err := tpl.ExecuteTemplate(w, template_file, data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}
