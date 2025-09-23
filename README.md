sudo apt update && apt upgrade -y
sudo apt install golang-go -y
go version

#bot start
go run cmd/main.go
