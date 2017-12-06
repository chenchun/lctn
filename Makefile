PACKAGE=github.com/chenchun/lctn

all:build
build:
	@mkdir -p bin
	@env GOOS=linux GOARCH=amd64 go build -v -o bin/lctn ${PACKAGE}/cmd/lctn
test:build
	@bats tests/test.bats
clean:
	@rm -rf bin
