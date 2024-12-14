#!/bin/bash
set -e
echo "Test Agent..."

echo "Preparing python dependency"
sudo pip3 install -r test/mock_driver/clock/requirements.txt

echo "Running Mock E2E test server for openapi"
sudo python3 test/mock_driver/clock/clock.py &

echo "Start Agent:"
go build -o cmd/agent/agent ./cmd/agent
cmd/agent/agent --env=test/mock_driver/clock/agent_env &

echo "Start Test Client:"
go build -o test/agent_test ./cmd/test/agent
cd test
./agent_test mini --agent-id clock-agent --has-inited
cd ..
