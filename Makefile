BINARY_NAME=./bin/main

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux

test:
	go test ./...

coverage:
	rm -r ./coverage && mkdir ./coverage
	go test ./... -coverprofile=./coverage/coverage.out && go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html

lint:
	golangci-lint run --fix -c ./.golangci.yml