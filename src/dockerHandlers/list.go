/*
author Vyacheslav Kryuchenko
*/
package dockerHandlers

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func List(connectedClient *client.Client, showAll bool) ([]types.Container, error) {
	return connectedClient.ContainerList(
		context.Background(),
		types.ContainerListOptions{All: showAll},
	)
}
