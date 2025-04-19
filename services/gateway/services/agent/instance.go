package agent

import "github.com/aenjoy/iot-lubricant/pkg/docker"

func bootAgentInstance(containerId string) error {
	return docker.StartContainer(containerId)
}
