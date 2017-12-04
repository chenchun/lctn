# lctn

lctn is a simple command line program to run a process in linux container. It focuses on a brief golang code description of how to implement a linux container runtime.

# Usage

Be sure to run as a root user since creating namespace requires SYS_ADMIN capability.

```
$ go get github.com/chenchun/lctn
$ cd $GOPATH/src/github.com/chenchun/lctn
$ make
github.com/chenchun/lctn
$ ./lctn -logtostderr -root `pwd`/rootfs /bin/sh
/ # env
SHLVL=1
PATH=/sbin:/bin:/usr/sbin:/usr/bin
PWD=/
/ # ls
bin      dev      etc      hello    lib      linuxrc  proc     sbin     sys      usr      var
/ # exit
```

# Note

`rootfs` directory in this repository comes from https://hub.docker.com/r/chenchun/hello/.
You can easily build your own rootfs from a docker image.

```
$ mkdir rootfs && cd rootfs
$ CID=$(docker run -d $image)
$ docker export $CID | tar xf -
```