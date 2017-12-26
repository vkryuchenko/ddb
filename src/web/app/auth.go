/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"net/http"
)

func (p *Provider) authPage(writer http.ResponseWriter, request *http.Request) {
	err := p.execTemplate(writer, "auth.tmpl", nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
