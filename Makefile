PACKAGE=github.com/chenchun/lctn

all:build
build:
	@mkdir -p bin
	@env GOOS=linux GOARCH=amd64 go build -v -o bin/lctn ${PACKAGE}/cmd/lctn
	@echo 'Successfully compile, try executing: bin/lctn -logtostderr -root `pwd`/rootfs /bin/sh'
test:build
	# This is needed on travis-ci, CAP_NET_BIND_SERVICE bind a socket to Internet
	# domain privileged ports (port numbers less than 1024)
	@sudo setcap CAP_NET_BIND_SERVICE+epi rootfs/hello
	@bats tests/test.bats
clean:
	@rm -rf bin
