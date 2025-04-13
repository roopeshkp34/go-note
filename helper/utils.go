package helper

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := filepath.Join("templates", tmpl)
	t := template.Must(template.ParseFiles(tmplPath))
	t.Execute(w, data)
}
