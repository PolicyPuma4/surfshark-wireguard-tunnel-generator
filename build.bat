set GOOS=windows
set GOARCH=amd64
go.exe build -o swtr_windows_amd64.exe -ldflags="-s -w" cmd\surfshark-wireguard-tunnel-generator\main.go

set GOOS=linux
set GOARCH=amd64
go.exe build -o swtr_linux_amd64 -ldflags="-s -w" cmd\surfshark-wireguard-tunnel-generator\main.go
