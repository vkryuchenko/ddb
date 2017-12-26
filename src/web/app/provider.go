/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"goji.io"
	"goji.io/pat"
	"html/template"
	"log"
	"net/http"
	"strings"
	"web"
)

type Provider struct {
	Listen          string
	ApplicationName string
	Secret          string
	instance        *goji.Mux
	templates       *template.Template
}

func (p *Provider) initTemplates() {
	resources := web.AssetNames()
	allTemplates := template.New("")
	for _, resName := range resources {
		if !strings.HasSuffix(resName, ".tmpl") {
			continue
		}
		bytes, err := web.Asset(resName)
		if err != nil {
			log.Fatal(err)
		}
		allTemplates, err = allTemplates.New(resName).Parse(string(bytes[:]))
	}
	p.templates = allTemplates
	//log.Printf(allTemplates.DefinedTemplates())
}

func (p *Provider) init() {
	mux := goji.NewMux()
	// list of handlers
	mux.HandleFunc(pat.Get("/static/*"), staticFile)
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
