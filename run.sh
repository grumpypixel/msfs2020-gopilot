#!/bin/bash

EXEC=./gopilot.exe
if [ ! -f "$EXEC" ]; then
  echo "$EXEC does not exists..."
  ./build.sh
fi

echo "starting..."
./$EXEC
