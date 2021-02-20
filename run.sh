#!/bin/bash

EXEC=./gopilot.exe
if [ ! -f "$EXEC" ]; then
  echo "$EXEC does not exists..."
  ./build.sh
fi

if [ "$1" = "clean" ]; then
  echo "Cleaning..."
  DLL=./SimConnect.dll
  if [ -f "$DLL" ]; then
    rm $DLL
  fi
fi

echo "Starting..."
# ./$EXEC --name COCO --searchpath "." --address 0.0.0.0:8888 --requestinterval 250 --timeout 600  --verbose "true"
./$EXEC
