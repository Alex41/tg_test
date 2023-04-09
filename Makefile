swagger:
	swag init --parseDependency \
		-g server/server.go \
		-o docs -ot json

build:
	go build -ldflags="-s -w" -o ./bin/app main/main.go

install_swagger:
	go get -u github.com/swaggo/swag/cmd/swag@v1.8.7
	go install github.com/swaggo/swag/cmd/swag@v1.8.7

go_download:
	go mod download

go_tidy:
	go mod tidy
