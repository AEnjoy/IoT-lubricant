package agent

import "github.com/AEnjoy/IoT-lubricant/pkg/docker"

func bootAgentInstance(containerId string) error {
	return docker.StartContainer(containerId)
}
