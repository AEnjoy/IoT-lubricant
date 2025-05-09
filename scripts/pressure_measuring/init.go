package main

import (
	"os"
	"time"
)

func init() {
	dataCh = make(chan *Data, maxBuffer)

	var ok bool
	userID, ok = os.LookupEnv(ENV_USER_ID)
	if !ok {
		panic("ENV_USER_ID not set")
	}
	gatewayID, ok = os.LookupEnv(ENV_GATEWAY_ID)
	if !ok {
		panic("ENV_GATEWAY_ID not set")
	}
	agentIDfile, ok = os.LookupEnv(ENV_AGENT_ID_FILES)
	if !ok {
		panic("ENV_AGENT_ID_FILES not set")
	}
	hostAddress, ok = os.LookupEnv(ENV_HOST_ADDRESS)
	if !ok {
		panic("ENV_HOST_ADDRESS not set")
	}
	algorithm = os.Getenv(ENV_COMPRESS_ALGORITHM)

	if err := initAgentIDPools(agentIDfile); err != nil {
		panic(err)
	}

	startTime = time.Now()
	dataTime = startTime.Add(-25 * time.Hour)
}
