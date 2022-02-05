.PHONY: frontend
all: frontend
	go build .

linux: frontend
	GOARCH=amd64 GOOS=linux go build .

test:
	go test ./...
	cd frontend && yarn lint

frontend:
	cd frontend && yarn && yarn build && cd ..
	cp -R frontend/dist cmd/
