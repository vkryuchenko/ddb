/*
author Vyacheslav Kryuchenko
*/
package app

import (
	"dockerHandlers"
	"github.com/docker/docker/api/types"
	"net/http"
	"os"
	"path/filepath"
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
	Username        string
	Quota           uint
	AvailableImages map[string]string
	Containers      []DockerContainer
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

func (uc *UserContainers) filter(containers []types.Container, images []string) error {
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
	uc.AvailableImages = make(map[string]string)
	if len(uc.Containers) < 1 || len(images) == 0 {
		for _, image := range images {
			uc.AvailableImages[image] = image
		}
	} else {
		for _, cnt := range uc.Containers {
			for _, image := range images {
				if image == cnt.Image {
					continue
				}
				uc.AvailableImages[image] = image
			}
		}
	}
	return nil
}

func (uc *UserContainers) collectInfo(p *Provider) error {
	containersList, err := dockerHandlers.List(p.Docker.Client, true)
	if err != nil {
		return err
	}
	return uc.filter(containersList, p.Docker.Images)
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

func (p *Provider) actionCreateContainer(w http.ResponseWriter, r *http.Request) {
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
	image := r.FormValue("image")
	containerName := strings.Replace(login, ".", "_", -1) +
		"_" +
		strings.Replace(
			strings.Replace(
				strings.Replace(
					image,
					":",
					"_",
					-1,
				),
				"/",
				"_",
				-1,
			),
			".",
			"_",
			-1,
		)
	containerData := filepath.Join(p.Docker.DataRoot, containerName)
	err = os.MkdirAll(containerData, 0776)
	if err != nil && err != os.ErrExist {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = dockerHandlers.CreateAndRun(p.Docker.Client, image, containerName, containerData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
