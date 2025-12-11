# Linux 版本
GOOS=linux GOARCH=amd64 go build -o build/cargo-m-linux-amd64-1.1.0 ./cmd/

# Windows 版本
GOOS=windows GOARCH=amd64 go build -o build/cargo-m-windows-amd64-1.1.0.exe ./cmd/