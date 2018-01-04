/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"net/http"
)

func (p *Provider) loginPage(w http.ResponseWriter, r *http.Request) {
	err := p.execTemplate(w, "login.tmpl", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *Provider) checkLogin(w http.ResponseWriter, r *http.Request) {}
