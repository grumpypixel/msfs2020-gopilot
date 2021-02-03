#!/bin/bash

EXEC=./gopilot.exe
if test -f "$EXEC"; then
  echo "removing old executable..."
  rm $EXEC
fi

echo "building..."
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o $EXEC gopilot/main.go gopilot/request_manager.go
echo "done"
