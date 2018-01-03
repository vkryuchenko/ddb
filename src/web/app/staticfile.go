/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"web"
)

func serveFile(resourceName string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	resourcePath := filepath.FromSlash(filepath.Join(wd, "src", "web", "resources", resourceName))
	file, err := os.Open(resourcePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func serveAsset(resourceName string) ([]byte, error) {
	return web.Asset(resourceName)
}

func (p *Provider) staticFile(writer http.ResponseWriter, request *http.Request) {
	var content []byte
	var serveErr error
	mimeType := "text/plain; charset=utf-8"
	resourceName := request.URL.Path[1:]
	if !p.Develop {
		content, serveErr = serveAsset(resourceName)
	} else {
		content, serveErr = serveFile(resourceName)
	}
	if serveErr != nil {
		http.Error(writer, serveErr.Error(), http.StatusNotFound)
		return
	}
	split := strings.Split(resourceName, ".")
	resourceType := split[len(split)-1]
	switch resourceType {
	case "css":
		mimeType = "text/css; charset=utf-8"
	case "js":
		mimeType = "text/javascript; charset=utf-8"
	case "svg":
		mimeType = "image/svg+xml; charset=utf-8"
	case "woff":
		mimeType = "application/font-woff; charset=utf-8"
	case "woff2":
		mimeType = "application/font-woff2; charset=utf-8"
	default:
		http.Error(writer, "invalid type", http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", mimeType)
	writer.Write(content)
}
