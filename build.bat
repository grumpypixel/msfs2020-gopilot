echo off
set EXEC=gopilot.exe
set TEMPLATE=./tools/packifier/template.gopher
set PACKAGE=app
set FUNCTION=SimConnectDLL
set DLL_IN="C:/MSFS SDK/SimConnect SDK/lib/SimConnect.dll"
set DLL_OUT="C:/MSFS SDK/SimConnect SDK/lib/SimConnect.dll"

if exist %EXEC% (
    echo removing old executable
    del %EXEC%
)

IF "%1"=="release" (
    echo pack simconnect.dll
    go run ./tools/packifier/main.go --in %DLL_IN% --out %DLL_OUT% --template %TEMPLATE% --package %PACKAGE% --function %FUNCTION%
)

echo building %EXEC%
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
go build -o %EXEC% ./cmd/main.go
echo done.
