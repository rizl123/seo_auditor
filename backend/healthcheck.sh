#!/bin/sh

if lsof -i :8080 -sTCP:LISTEN > /dev/null; then
    exit 0
else
    echo "Healthcheck failed"
    exit 1
fi
