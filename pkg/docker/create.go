package docker

import (
	"context"
	"net/http"

	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func Create(ctx context.Context, c types.Container) (*container.CreateResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	// Pull image if necessary
	switch c.Source.PullWay {
	case types.ImagePullFromBinary:
		if err := pullFromImageBinaryData(ctx, cli, c.Source.FromBinary); err != nil {
			return nil, err
		}
	case types.ImagePullFromUrl:
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
	case types.ImagePullFromRegistry:
		if _, err := cli.ImagePull(ctx, c.Source.FromRegistry, image.PullOptions{RegistryAuth: c.Source.RegistryAuth}); err != nil {
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
			c.Network: {},
		},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, netConfig, nil, c.Name)
	if err != nil {
		return nil, err
	}

	return &resp, cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
}
