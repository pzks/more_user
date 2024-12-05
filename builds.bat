@echo off
echo start build inux amd64
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ./bin/linux_amd64/mu ./cmd/more_user

echo start build windows amd64
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o ./bin/windows_amd64/mu.exe ./cmd/more_user