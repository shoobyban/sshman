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
	rm -f /Users/sam/proj/sshman/testdata
	rm -f ./sshman.port
	rm -f ./test.log
	SSHMAN_STORAGE=/Users/sam/proj/sshman/testdata ./sshman_test web | tee -a test.log &
	sleep 1
	cat sshman.port
	cd frontend && yarn wdio --baseUrl=http://localhost:$$(cat ../sshman.port) || echo "failure!"
	killall sshman_test
	rm -f ./sshman_test
	open frontend/test/reports/suite-0-0/0-0/report.html
	open frontend/test/reports/suite-0-0/0-1/report.html

frontend:
	cd frontend && yarn && yarn build && cd ..
	cp -R frontend/dist cmd/
