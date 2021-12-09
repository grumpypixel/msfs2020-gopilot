#!/bin/bash

EXEC=./gopilot.exe

# rm $EXEC

if [ ! -f "$EXEC" ]; then
  echo "$EXEC does not exists..."
  ./scripts/build.sh
fi

# if [ "$1" = "clean" ]; then
#   echo "Cleaning..."
#   DLL=./SimConnect.dll
#   if [ -f "$DLL" ]; then
#     rm $DLL
#   fi
# fi

echo "Starting..."
./$EXEC --cfg configs/config.yml
