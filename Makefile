all:
	cd ui && yarn && yarn build
	cp -R ui/dist cmd/
	go build .

linux:
	GOARCH=amd64 GOOS=linux go build .

test:
	go test ./...
