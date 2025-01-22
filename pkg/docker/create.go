package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	docker "github.com/AEnjoy/IoT-lubricant/pkg/types/container"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func UpdateContainer(ctx context.Context, c *docker.Container, oldName string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	if err = pullImage(ctx, cli, c); err != nil {
		return "", err
	}
	if err = cli.ContainerStop(ctx, oldName, container.StopOptions{}); err != nil {
		logger.Warnf("Failed to stop container %s: %v", oldName, err)
	}
	if err = cli.ContainerRemove(ctx, oldName, container.RemoveOptions{Force: true}); err != nil {
		return "", fmt.Errorf("failed to remove container %s: %w", oldName, err)
	}

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
		Mounts:       c.Mount,
	}
	netConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			c.Network: {},
		},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, netConfig, nil, c.Name)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return resp.ID, fmt.Errorf("failed to start container: %w", err)
	}

	logger.Infof("Container %s updated successfully with new image %s", c.Name, c.Source.RegistryPath)
	return resp.ID, nil
}

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
	if err = pullImage(ctx, cli, c); err != nil {
		return nil, err
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
