/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"html/template"
	"net/http"
)

func (p *Provider) mainPage(writer http.ResponseWriter, _ *http.Request) {
	mainPage(writer, p.templates)
}

func mainPage(writer http.ResponseWriter, templ *template.Template) {
	err := templ.ExecuteTemplate(writer, "main.tmpl", nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
