FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o oim ./cmd/api/*.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/oim .

COPY --from=builder /app/config.json .

EXPOSE 8080

CMD ["./oim"]