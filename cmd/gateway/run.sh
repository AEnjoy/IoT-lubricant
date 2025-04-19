#!/bin/bash

go build -v

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

if [ -f "$SCRIPT_DIR/gateway.env" ]; then
    export $(grep -v '^#' "$SCRIPT_DIR/gateway.env" | xargs)
else
    echo "Error: $SCRIPT_DIR/gateway.env does not exist." >&2
    exit 1
fi

./gateway
