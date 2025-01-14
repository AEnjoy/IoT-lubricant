package docker

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	docker "github.com/AEnjoy/IoT-lubricant/pkg/types/container"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func Create(ctx context.Context, c *docker.Container) (*container.CreateResponse, error) {
	if c.Compose != nil {
		err := os.WriteFile("docker-compose.yaml", *c.Compose, 0644)
		if err != nil {
			logger.Error("write docker-compose.yaml failed", err)
			return nil, err
		}
		err = exec.Command("docker-compose", "up", "-d").Run()
		if err != nil {
			logger.Error("docker-compose up failed", err)
			return nil, err
		}
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	createNetwork(cli)

	// Pull image if necessary
	switch c.Source.PullWay {
	case docker.ImagePullFromBinary:
		if err := pullFromImageBinaryData(ctx, cli, c.Source.FromBinary); err != nil {
			return nil, err
		}
	case docker.ImagePullFromUrl:
		pullReq, err := http.NewRequest(http.MethodGet, c.Source.FromUrl, nil)
		if err != nil {
			return nil, err
		}
		resp, err := cli.HTTPClient().Do(pullReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if err := pullFromImageBinaryReader(ctx, cli, resp.Body); err != nil {
			return nil, err
		}
	case docker.ImagePullFromRegistry:
		path := func() string {
			c.Source.FromRegistry = strings.Trim(c.Source.FromRegistry, "/")
			c.Source.FromRegistry = strings.TrimLeft(c.Source.FromRegistry, "/")
			c.Source.RegistryPath = strings.Trim(c.Source.RegistryPath, "/")
			c.Source.RegistryPath = strings.TrimLeft(c.Source.RegistryPath, "/")
			if len(c.Source.FromRegistry) == 0 {
				c.Source.FromRegistry = "docker.io"
			}
			return c.Source.FromRegistry + "/" + c.Source.RegistryPath
		}()
		if _, err := cli.ImagePull(ctx, path, image.PullOptions{RegistryAuth: c.Source.RegistryAuth}); err != nil {
			return nil, err
		}
	default:
		return nil, ErrPullWay
	}

	// configs
	config := &container.Config{
		Image:        c.Source.RegistryPath,
		Env:          convertEnvMap(c.Env),
		ExposedPorts: convertExposePort(c.ExposePort),
	}
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyAlways,
		},
		PortBindings: convertPortBindings(c.ExposePort),
	}
	netConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			c.Network:   {},
			NetWorkName: {NetworkID: NetWorkName},
		},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, netConfig, nil, c.Name)
	if err != nil {
		return nil, err
	}

	return &resp, cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
}

// DeployAgent 部署agent容器 返回 agent-container id
func DeployAgent() (string, error) {
	resp, err := Create(context.Background(), &model.AgentContainer)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}
