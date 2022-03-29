BINARY_NAME=./bin/gomstore

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} .

clean:
	go clean
	rm ${BINARY_NAME}

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./... -c ./.golangci.yml