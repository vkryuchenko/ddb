package dockerHandlers

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func CreateAndRun(connectedClient *client.Client, image string, name string, mountPath string) error {
	ctx := context.Background()
	_, err := connectedClient.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	cntCfg := container.Config{
		Image: image,
	}
	portBind := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "",
	}
	portMap := nat.PortMap{
		"5432/tcp": []nat.PortBinding{
			portBind,
		},
	}
	hostCfg := container.HostConfig{
		Binds:        []string{mountPath + ":/var/lib/postgresql/data"},
		AutoRemove:   true,
		PortBindings: portMap,
	}
	result, err := connectedClient.ContainerCreate(ctx, &cntCfg, &hostCfg, nil, name)
	if err != nil {
		return err
	}
	return connectedClient.ContainerStart(ctx, result.ID, types.ContainerStartOptions{})
}
