package docker

import (
	"context"

	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

const (
	DefaultSubnet  = "172.29.26.0/20"
	DefaultGateway = "172.29.26.1"
)

// todo: 使用本机唯一标识 作为网桥名称 如主板序号
var NetWorkName = "lubricant_edge_network"

func createNetwork(cli *client.Client) {
	_, _ = cli.NetworkCreate(context.Background(), NetWorkName, network.CreateOptions{
		Driver: network.NetworkBridge,
		IPAM: &network.IPAM{
			Driver: network.NetworkDefault,
			Config: []network.IPAMConfig{
				{
					Subnet:  DefaultSubnet,
					Gateway: DefaultGateway,
				},
			},
		},
	})
}
