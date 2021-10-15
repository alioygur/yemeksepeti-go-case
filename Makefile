build-docker: test
	docker build -t herodb .
build:
	go build -o restapi ./cmd/restapi
test:
	go test -v ./...
run: test build
	PORT=3000 ./restapi
run-dev: test
	cd ./cmd/restapi && go run .
	