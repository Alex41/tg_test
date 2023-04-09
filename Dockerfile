FROM golang:1.19-buster AS builder
WORKDIR /app

COPY Makefile .
COPY go.mod .
COPY go.sum .

RUN make install_swagger

COPY . .

RUN make go_download go_tidy swagger build

FROM buildpack-deps:18.04-scm
WORKDIR /root

COPY --from=builder /app/bin .
COPY --from=builder /app/docs ./docs
COPY ./migrations ./migrations

ENTRYPOINT ["./app"]
