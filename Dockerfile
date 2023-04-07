FROM golang:1.19-alpine3.15 AS builder
WORKDIR /app

COPY . .

RUN go build -ldflags="-s -w" -o ./bin/app main/main.go

FROM alpine:3.15
WORKDIR /root

COPY --from=builder /app/bin .
COPY ./migrations ./migrations

RUN chmod +x ./app

ENTRYPOINT ["./app"]
