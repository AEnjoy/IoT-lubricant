package main

const (
	maxBuffer = 200000
)

const (
	ENV_USER_ID        = "ENV_USER_ID"
	ENV_GATEWAY_ID     = "ENV_GATEWAY_ID"
	ENV_AGENT_ID_FILES = "ENV_AGENT_ID_FILES"

	ENV_COMPRESS_ALGORITHM = "ENV_COMPRESS_ALGORITHM"
	ENV_HOST_ADDRESS       = "ENV_HOST_ADDRESS"
)

var (
	userID      string
	gatewayID   string
	agentIDfile string

	algorithm string

	hostAddress string
)
