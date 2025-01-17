package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	docker "github.com/AEnjoy/IoT-lubricant/pkg/types/container"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Deprecated: 该接口应该直接调用，不应该通过该接口来调用
type Api struct {
	ctx context.Context
	cli *client.Client
}

// Deprecated: 该接口应该直接调用，不应该通过该接口来调用
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

// Deprecated: use New Docker Api
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

func convertEnvMap(env map[string]string) []string {
	var envList []string
	for key, value := range env {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}
	return envList
}
func convertPortBindings(port map[string]int) nat.PortMap {
	portBindings := make(nat.PortMap)
	for containerPort, hostPort := range port {
		portBindings[nat.Port(containerPort)] = []nat.PortBinding{
			{
				HostPort: fmt.Sprintf("%d", hostPort),
			},
		}
	}
	return portBindings
}

func convertExposePort(exposePort map[string]int) nat.PortSet {
	portSet := make(nat.PortSet)
	for containerPort, _ := range exposePort {
		portSet[nat.Port(containerPort)] = struct{}{}
	}
	return portSet
}

func pullFromImageBinaryReader(ctx context.Context, cli *client.Client, data io.Reader) error {
	_, err := cli.ImageLoad(ctx, data, false)
	if err != nil {
		return err
	}
	return nil
}
func pullFromImageBinaryData(ctx context.Context, cli *client.Client, data []byte) error {
	return pullFromImageBinaryReader(ctx, cli, bytes.NewReader(data))
}
func pullImage(ctx context.Context, cli *client.Client, c *docker.Container) error {
	switch c.Source.PullWay {
	case docker.ImagePullFromBinary:
		if err := pullFromImageBinaryData(ctx, cli, c.Source.FromBinary); err != nil {
			return err
		}
	case docker.ImagePullFromUrl:
		pullReq, err := http.NewRequest(http.MethodGet, c.Source.FromUrl, nil)
		if err != nil {
			return err
		}
		resp, err := cli.HTTPClient().Do(pullReq)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if err := pullFromImageBinaryReader(ctx, cli, resp.Body); err != nil {
			return err
		}
	case docker.ImagePullFromRegistry:
		path := func() string {
			c.Source.FromRegistry = strings.Trim(c.Source.FromRegistry, "/")
			c.Source.RegistryPath = strings.Trim(c.Source.RegistryPath, "/")
			if len(c.Source.FromRegistry) == 0 {
				c.Source.FromRegistry = "docker.io"
			}
			return c.Source.FromRegistry + "/" + c.Source.RegistryPath
		}()
		if _, err := cli.ImagePull(ctx, path, image.PullOptions{RegistryAuth: c.Source.RegistryAuth}); err != nil {
			return err
		}
	default:
		return ErrPullWay
	}
	return nil
}
