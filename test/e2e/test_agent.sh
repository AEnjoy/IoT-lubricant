#!/bin/bash
set -e
echo "Test Agent..."

echo "Preparing python dependency"
sudo pip3 install -r test/mock_driver/clock/requirements.txt

echo "Running Mock E2E test server for openapi"
python3 test/mock_driver/clock/clock.py &

echo "Start Agent:"
cd cmd/agent
go build . -o agent
./agent --env=core.env &

echo "Start Test Client:"
cd ../test/agent
go build . -o agent_test
cd ../../test
../cmd/test/agent/agent_test mini --agent-id clock-agent

