.PHONY: all build backend frontend linux test check clean

all: build

# build is the default target, builds for the current OS/ARCH
build: frontend
	go build -o sshman .

# backend is an alias for build
backend: build

# frontend builds the web UI
frontend:
	cd frontend && npm install && npm run build
	rm -rf cmd/dist
	cp -R frontend/dist cmd/

# linux target for cross-compilation
linux: frontend
	GOOS=linux GOARCH=amd64 go build -o sshman-linux .

# test runs backend and frontend tests
test: frontend
	go test ./...
	cd frontend && npm run lint

# check is an alias for test, often used in CI
check: test

# clean removes build artifacts
clean:
	rm -f sshman sshman-linux
	rm -rf cmd/dist
	cd frontend && rm -rf dist node_modules

