/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"html/template"
	"net/http"
)

func (p *Provider) authPage(writer http.ResponseWriter, request *http.Request) {
	authPage(writer, p.templates)
}

func authPage(writer http.ResponseWriter, templ *template.Template) {
	err := templ.ExecuteTemplate(writer, "auth.tmpl", nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
