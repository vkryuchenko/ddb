package app

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"web"
)

func (p *Provider) execTemplate(writer http.ResponseWriter, templateName string, data interface{}) error {
	var tmpl *template.Template
	if !p.Develop {
		tmpl = p.templates
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		tmplPattern := filepath.Join(wd, "src", "web", "resources", "*.tmpl")
		t, err := template.ParseGlob(tmplPattern)
		if err != nil {
			return err
		}
		tmpl = t
	}
	return tmpl.ExecuteTemplate(writer, templateName, data)
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
}
