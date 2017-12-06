PACKAGE=github.com/chenchun/lctn

all:build
build:
	@mkdir -p bin
	@env GOOS=linux GOARCH=amd64 go build -v -o bin/lctn ${PACKAGE}/cmd/lctn
clean:
	@rm -rf bin
