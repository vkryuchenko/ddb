/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"goji.io"
	"goji.io/pat"
	"helpers"
	"html/template"
	"log"
	"net/http"
)

type Provider struct {
	Listen          string
	ApplicationName string
	Secret          string
	Develop         bool
	LDAPClient      *helpers.LDAPClient
	EmailClient     *helpers.EmailClient
	instance        *goji.Mux
	templates       *template.Template
}

func (p *Provider) init() {
	mux := goji.NewMux()
	// list of handlers
	mux.HandleFunc(pat.Get("/static/*"), p.staticFile)
	mux.HandleFunc(pat.Get("/"), p.mainPage)
	mux.HandleFunc(pat.Get("/auth"), p.authPage)
	// end list of handlers
	p.instance = mux
}

func StartServer(p *Provider) {
	p.initTemplates()
	p.init()
	log.Fatal(http.ListenAndServe(p.Listen, p.instance))
}
