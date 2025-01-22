package docker

import (
	"context"

	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func StartContainer(id string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Error("docker api init failed", err)
		return err
	}
	defer cli.Close()

	return cli.ContainerStart(context.Background(), id, container.StartOptions{})
}
func StopContainer(id string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Error("docker api init failed", err)
		return err
	}
	defer cli.Close()
	return cli.ContainerStop(context.Background(), id, container.StopOptions{})
}

func IsContainerRunning(id string) bool {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Error("docker api init failed", err)
		return false
	}
	defer cli.Close()

	inspect, err := cli.ContainerInspect(context.Background(), id)
	if err != nil {
		logger.Error("inspect container failed", err)
		return false
	}
	return inspect.State.Running
}
