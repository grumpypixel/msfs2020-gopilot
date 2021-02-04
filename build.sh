#!/bin/bash

EXEC=./gopilot.exe
if test -f "$EXEC"; then
  echo "Removing old executable..."
  rm $EXEC
fi

if [ "$1" = "release" ]; then
  TEMPLATE="./tools/packifier/template.gopher"
  PACKAGE="main"

  # pack dll
  go run ./tools/packifier/main.go --in "/c/MSFS SDK/Samples/SimvarWatcher/bin/x64/Release/SimConnect.dll" --out "./gopilot/dllpack.go" --template $TEMPLATE --package $PACKAGE --function PackedSimConnectDLL

  # pack assets
  ASSETS_TAR="assets.tar"
  go run ./tools/tarifier/main.go --in "assets" --out ""
  go run ./tools/packifier/main.go --in $ASSETS_TAR --out "./gopilot/assetspack.go" --template $TEMPLATE --package $PACKAGE --function PackedAssets
  rm $ASSETS_TAR
fi

echo "Building..."
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o $EXEC gopilot/main.go gopilot/request_manager.go gopilot/assetspack.go gopilot/dllpack.go
echo "Done."
