.PHONY: frontend
all: backend frontend

backend: frontend
	go build .

linux: frontend
	GOARCH=amd64 GOOS=linux go build .

test: backend frontend
	go test ./...
	cd frontend && yarn lint

frontend:
	cd frontend && yarn && yarn build && cd ..
	cp -R frontend/dist cmd/
