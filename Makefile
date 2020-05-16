.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/postman postman/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/mailtruck mailtruck/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/mailman mailman/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
