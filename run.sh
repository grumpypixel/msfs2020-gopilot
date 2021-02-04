#!/bin/bash

EXEC=./gopilot.exe
if [ ! -f "$EXEC" ]; then
  echo "$EXEC does not exists..."
  ./build.sh
fi

#rm SimConnect.dll

echo "Starting..."
./$EXEC
