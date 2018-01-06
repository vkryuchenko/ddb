/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"dockerHandlers"
	"github.com/docker/docker/api/types"
	"net/http"
	"strings"
)

const (
	queryGetLoginBySessionId = `
SELECT login
FROM public.users
  JOIN public.sessions ON users.id = sessions.user_id
WHERE session_id = $1;
`
)

type DockerContainer struct {
	ID        string
	Name      string
	Image     string
	Port      uint16
	State     string
	DiskUsage uint
}

type UserContainers struct {
	Username   string
	Quota      uint
	Containers []DockerContainer
}

func (p *Provider) getLoginBySessionId(uuid string) (string, error) {
	login := ""
	db, err := p.Database.Connect()
	if err != nil {
		return login, err
	}
	query, err := db.Prepare(queryGetLoginBySessionId)
	if err != nil {
		return login, err
	}
	err = query.QueryRow(uuid).Scan(&login)
	if err != nil {
		return login, err
	}
	return login, nil
}

func (uc *UserContainers) filter(containers []types.Container) error {
	if len(containers) < 1 {
		return nil
	}
	prefix := "/" + strings.Replace(uc.Username, ".", "_", -1)
	for _, cnt := range containers {
		if strings.HasPrefix(cnt.Names[0], prefix) {
			dc := DockerContainer{
				ID:    cnt.ID,
				Name:  cnt.Names[0][1:],
				Image: cnt.Image,
				State: cnt.State,
			}
			if len(cnt.Ports) > 0 {
				dc.Port = cnt.Ports[0].PublicPort
			}
			uc.Containers = append(uc.Containers, dc)
		}
	}
	return nil
}

func (uc *UserContainers) collectInfo(p *Provider) error {
	containersList, err := dockerHandlers.List(p.Docker.Client, true)
	if err != nil {
		return err
	}
	return uc.filter(containersList)
}

func (p *Provider) mainPage(w http.ResponseWriter, r *http.Request) {
	ok := p.sessionValid(w, r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	cookie, err := r.Cookie(p.ApplicationName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	login, err := p.getLoginBySessionId(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := UserContainers{
		Username: login,
		Quota:    p.Docker.Quota / (1 << 30),
	}
	err = data.collectInfo(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = p.execTemplate(w, "main.tmpl", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
