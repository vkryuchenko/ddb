/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"net/http"
	"strings"
	"web"
)

func staticFile(writer http.ResponseWriter, request *http.Request) {
	mimeType := "text/plain; charset=utf-8"
	resourceName := request.URL.Path[1:]
	content, err := web.Asset(resourceName)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	split := strings.Split(resourceName, ".")
	resourceType := split[len(split)-1]
	switch resourceType {
	case "css":
		mimeType = "text/css; charset=utf-8"
	case "js":
		mimeType = "text/javascript; charset=utf-8"
	default:
		http.Error(writer, "invalid type", http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", mimeType)
	writer.Write(content)
}
