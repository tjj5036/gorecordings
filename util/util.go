package util

import (
	"html/template"
	"log"
	"net/http"
)

// AddIndexFunc allows you to take an index (or any number really) and add
// another value to it while rendering a template. Useful for setlists because
// index values start at 0 and that doesn't make sense from a setlist context
func addIndexFunc(index int, increment int) int {
	return index + increment
}

// RenderTemplate renders a template given a template file to render and data
// in the form of an interface
func RenderTemplate(w http.ResponseWriter, template_file string, data interface{}) {
	funcMap := template.FuncMap{
		"add": addIndexFunc,
	}
	var tpl *template.Template
	tpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*html"))
	err := tpl.ExecuteTemplate(w, template_file, data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}
