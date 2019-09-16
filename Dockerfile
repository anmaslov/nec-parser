FROM golang:latest as builder

LABEL maintainer="Aleksey Maslov <a@anmaslov.ru>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nec-parser .

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/nec-parser .

CMD ["./nec-parser"]