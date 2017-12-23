PACKAGE=github.com/chenchun/lctn

all:build
build:
	@mkdir -p bin
	@env GOOS=linux GOARCH=amd64 go build -v -o bin/lctn ${PACKAGE}/cmd/lctn
	# lctn needs CAP_SYS_ADMIN to create namespaces, CAP_SYS_CHROOT to chroot,
	# CAP_SETUID/CAP_SETGID to mapping uid/gid of the new userspace.
	# Please refer to http://man7.org/linux/man-pages/man7/capabilities.7.html
	@sudo setcap CAP_SYS_ADMIN,CAP_SYS_CHROOT,CAP_SETUID,CAP_SETGID+epi bin/lctn
test:build
	# This is needed on travis-ci, CAP_NET_BIND_SERVICE bind a socket to Internet
	# domain privileged ports (port numbers less than 1024)
	@sudo setcap CAP_NET_BIND_SERVICE+epi rootfs/hello
	@bats tests/test.bats
clean:
	@rm -rf bin
