package app

import (
	"net/http"
	"strings"
)

func (p *Provider) cookieMiddleware(http.Handler) http.Handler {
	middleware := func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/login" {
			p.loginPage(writer, request)
			return
		}
		if strings.HasPrefix(request.URL.Path, "/static/") {
			p.staticFile(writer, request)
			return
		}
		_, err := request.Cookie(p.ApplicationName)
		if err != nil {
			http.Redirect(writer, request, "/login", http.StatusFound)
			return
		}
	}
	return http.HandlerFunc(middleware)
}

func (p *Provider) dropCookie(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
