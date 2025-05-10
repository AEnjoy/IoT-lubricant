#!/usr/bin/env bash

CGO_ENABLE=0 go build -v -o ./pressure-measuring .

export ENV_USER_ID="5a4ad48c-4985-4c12-b3de-d25c64427fdf"
export ENV_GATEWAY_ID="lubricant-gateway-0"
export ENV_AGENT_ID_FILES="agent_id.txt"
export ENV_HOST_ADDRESS="127.0.0.1:5423" #lubricant-grpcserver.lubricant.svc.cluster.local:5423
export ENV_COMPRESS_ALGORITHM="default"
./pressure-measuring
