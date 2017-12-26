/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"net/http"
)

func (p *Provider) mainPage(writer http.ResponseWriter, request *http.Request) {
	_, err := request.Cookie(p.ApplicationName)
	if err != nil {
		http.Redirect(writer, request, "/auth", http.StatusFound)
		return
	}
	err = p.execTemplate(writer, "main.tmpl", nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
