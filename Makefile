all:build
build:
	@env GOOS=linux GOARCH=amd64 go build -v -o lctn .
clean:
	@rm -f lctn
