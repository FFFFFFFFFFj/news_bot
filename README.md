sudo apt update && apt upgrade -y
sudo apt install golang-go -y
go version


#conection psql
psql "postgres://news_user:supersecret@127.0.0.1:5433/news_bot"


#bot start
go run cmd/main.go
