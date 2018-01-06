package helpers

import (
	"dockerHandlers"
	"github.com/docker/docker/client"
	"log"
)

type DockerConfig struct {
	DataRoot   string   `json:"data_root"`
	Quota      uint     `json:"quota"`
	Target     string   `json:"target"`
	Apiversion string   `json:"apiversion"`
	Images     []string `json:"images"`
	Client     *client.Client
}

func (dc *DockerConfig) InitClient() {
	var err error
	dc.Client, err = dockerHandlers.CreateClientConnection(dc.Target, dc.Apiversion)
	if err != nil {
		log.Panic(err)
	}
}
