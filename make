set GOARCH=amd64
set GOOS=linux
go build main

docker-compose up --build -d