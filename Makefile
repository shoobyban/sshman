.PHONY: frontend
all: backend frontend

backend: frontend
	go build .

linux: frontend
	GOARCH=amd64 GOOS=linux go build .

test: 
	go test ./...
	cd frontend && yarn && yarn lint

fulltest: backend frontend test
	mv sshman sshman_test
	rm -f ./testdata
	rm -f ./sshman.port
	SSHMAN_STORAGE=./testdata ./sshman_test web &
	sleep 1
	cat sshman.port
	cd frontend && yarn wdio --config ./wdio.conf.js  --baseUrl=http://localhost:$$(cat ../sshman.port)
	rm -f ./sshman_test

frontend:
	cd frontend && yarn && yarn build && cd ..
	cp -R frontend/dist cmd/
