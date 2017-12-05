PACKAGE=github.com/chenchun/lctn

all:build
build:
	@mkdir -p bin
	@env GOOS=linux GOARCH=amd64 go build -v -o bin/lctn ${PACKAGE}/cmd/lctn
	@env GOOS=linux GOARCH=amd64 go build -v -o bin/init ${PACKAGE}/cmd/init
clean:
	@rm -rf bin
