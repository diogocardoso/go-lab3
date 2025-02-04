FROM golang:1.21.3

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
COPY cmd/auction/.env /app/cmd/auction/.env

RUN go build -o /app/auction cmd/auction/main.go
EXPOSE 8080

ENTRYPOINT ["/app/auction"]