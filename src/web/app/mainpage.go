/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"html/template"
	"net/http"
)

func (p *Provider) mainPage(writer http.ResponseWriter, request *http.Request) {
	_, err := request.Cookie(p.ApplicationName)
	if err != nil {
		http.Redirect(writer, request, "/auth", http.StatusFound)
		return
	}
	mainPage(writer, p.templates)
}

func mainPage(writer http.ResponseWriter, templ *template.Template) {
	err := templ.ExecuteTemplate(writer, "main.tmpl", nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
