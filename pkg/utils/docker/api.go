package docker

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var (
	ErrNotInit = errors.New("docker api not init")
)

type Api struct {
	ctx context.Context
	cli *client.Client
}
type BinaryImageInfo struct {
	Reader    io.Reader // Tarball reader []byte
	ImageName string
}

func (a *Api) InstallFromImageBinary(imageData BinaryImageInfo, name string, exposePort int, env []string) error {
	if a.cli == nil {
		return ErrNotInit
	}

	_, err := a.cli.ImageLoad(a.ctx, imageData.Reader, false)
	if err != nil {
		return err
	}
	return a.Install(imageData.ImageName, name, exposePort, env)
}
func (a *Api) Install(imageLink, name string, exposePort int, env []string) error {
	if a.cli == nil {
		return ErrNotInit
	}
	//
	containerConfig := &container.Config{
		Image: imageLink,
		Env:   env,
		ExposedPorts: map[nat.Port]struct{}{
			nat.Port(fmt.Sprintf("%d/tcp", exposePort)): {},
		},
	}
	hostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			nat.Port(fmt.Sprintf("%d/tcp", exposePort)): {
				{
					HostPort: fmt.Sprintf("%d", exposePort),
				},
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}

	_, err := a.cli.ContainerCreate(a.ctx, containerConfig, hostConfig, nil, nil, name)
	if err != nil {
		return err
	}
	return a.cli.ContainerStart(a.ctx, name, container.StartOptions{})
}

func (a *Api) Remove(name string) error {
	return a.cli.ContainerRemove(a.ctx, name, container.RemoveOptions{})
}

func NewDockerApi(ctx context.Context) (*Api, error) {
	var dcli Api
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	dcli.ctx = ctx
	dcli.cli = cli
	return &dcli, nil
}
