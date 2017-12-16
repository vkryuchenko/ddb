/*
author Vyacheslav Kryuchenko
*/
package dockerHandlers

import (
	"github.com/docker/docker/client"
	"net/http"
)

func CreateClientConnection(address string, version string) (*client.Client, error) {
	httpClient := http.Client{Transport: http.DefaultTransport}
	httpHeaders := make(map[string]string)
	return client.NewClient(
		address,
		version,
		&httpClient,
		httpHeaders,
	)
}
