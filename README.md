# Lctn

[![Build Status](https://travis-ci.org/chenchun/lctn.svg?branch=master)](https://travis-ci.org/chenchun/lctn)
[![Go Report Card](https://goreportcard.com/badge/github.com/chenchun/lctn)](https://goreportcard.com/report/github.com/chenchun/lctn)

lctn is a simple command line program to run a process in a linux container. It focuses on a brief golang code description of how to implement a linux container runtime.

# Usage

Be sure to run as a root user since creating namespace requires SYS_ADMIN capability.

```
$ go get github.com/chenchun/lctn
$ cd $GOPATH/src/github.com/chenchun/lctn
$ make
github.com/chenchun/lctn
# run the following command to exec a shell inside a container
$ bin/lctn -logtostderr -root `pwd`/rootfs /bin/sh
/ # env
SHLVL=1
PATH=/sbin:/bin:/usr/sbin:/usr/bin
PWD=/
/ # ls
bin      dev      etc      hello    lib      linuxrc  proc     sbin     sys      usr      var
/ # exit
```

```
Usage of ./lctn [Flags] Command [Argument ...]:
  -alsologtostderr
    	log to standard error as well as files
  -init
    	is it the init process of container. (This flag is used by lctn internally)
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -root string
    	the root directory of container
  -stderrthreshold value
    	logs at or above this threshold go to stderr
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging

Command:
  The init command of container

Argument:
  The Arguments of container

Examples:
  lctn -logtostderr -root ./rootfs /bin/sh
```

# Note

`rootfs` directory in this repository comes from https://hub.docker.com/r/chenchun/hello/.
You can easily build your own rootfs from a docker image.

```
$ mkdir rootfs && cd rootfs
$ CID=$(docker run -d $image)
$ docker export $CID | tar xf -
```
