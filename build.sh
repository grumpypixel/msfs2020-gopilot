#!/bin/bash

EXEC=./gopilot.exe
if test -f "$EXEC"; then
  echo "Removing old executable..."
  rm $EXEC
fi

if [ "$1" = "release" ]; then
  TEMPLATE="./tools/packifier/template.gopher"
  PACKAGE="main"

  # # pack assets
  # ASSETS_DIR="assets"
  # ASSETS_TAR="$ASSETS_DIR.tar"
  # go run ./tools/tarifier/main.go --in $ASSETS_DIR --out ""
  # go run ./tools/packifier/main.go --in $ASSETS_TAR --out "./gopilot/assetspack.go" --template $TEMPLATE --package $PACKAGE --function PackedAssets
  # rm $ASSETS_TAR

  # # pack data
  # DATA_DIR="data"
  # DATA_TAR="$DATA_DIR.tar"
  # go run ./tools/tarifier/main.go --in $DATA_DIR --out ""
  # go run ./tools/packifier/main.go --in $DATA_TAR --out "./gopilot/datapack.go" --template $TEMPLATE --package $PACKAGE --function PackedData
  # rm $DATA_TAR

  # pack dll
  go run ./tools/packifier/main.go --in "/c/MSFS SDK/Samples/SimvarWatcher/bin/x64/Release/SimConnect.dll" --out "./gopilot/dllpack.go" --template $TEMPLATE --package $PACKAGE --function PackedSimConnectDLL
fi

echo "Building..."
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o $EXEC cmd/main.go
echo "Done."
