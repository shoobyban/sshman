all:
	go build .

linux:
	GOARCH=amd64 GOOS=linux go build .

test:
	go test ./...
