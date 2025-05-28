BINARY_NAME=golang-winservice.exe

build:
	go mod tidy
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME} ./cmd/service/main.go

test:
	go test ./...

clean:
	go clean
	rm ${BINARY_NAME}