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

func (p *Provider) mainPage(writer http.ResponseWriter, request *http.Request) {
	data := UserContainers{
		Username: "newton",
		Quota:    p.Docker.Quota / (1 << 30),
	}
	err := data.collectInfo(p)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = p.execTemplate(writer, "main.tmpl", data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
