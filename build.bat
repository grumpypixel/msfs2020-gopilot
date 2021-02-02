echo off
set EXEC=gopilot.exe

if exist %EXEC% (
    echo removing old executable
    del %EXEC%
)

echo building %EXEC%
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
go build -o %EXEC% gopilot/main.go gopilot/request_manager.go
echo done.
