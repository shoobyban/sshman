.PHONY: frontend
all: backend frontend

backend: frontend
	go build .

linux: frontend
	GOARCH=amd64 GOOS=linux go build .

test: backend frontend
	go test ./...
	cd frontend && npm run lint

frontend:
	cd frontend && npm i && npm run build && cd ..
	cp -R frontend/dist cmd/
