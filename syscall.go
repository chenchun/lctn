package main

import (
	"fmt"
	"runtime"
	"syscall"

	"github.com/golang/glog"
)

// Setns sets namespace using syscall. Note that this should be a method
// in syscall but it has not been added.
func Setns(ns int, nstype int) (err error) {
	_, _, e1 := syscall.Syscall(SYS_SETNS, uintptr(ns), uintptr(nstype), 0)
	if e1 != 0 {
		err = e1
	}
	return
}

// SYS_SETNS syscall allows changing the namespace of the current process.
var SYS_SETNS = map[string]uintptr{
	"386":     346,
	"amd64":   308,
	"arm64":   268,
	"arm":     375,
	"mips":    4344,
	"mipsle":  4344,
	"ppc64":   350,
	"ppc64le": 350,
	"s390x":   339,
}[runtime.GOARCH]

// GetNSfd gets a handle to a namespace
// identified by the path
func GetNSfd(path string) (int, error) {
	fd, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil {
		return -1, err
	}
	return fd, nil
}

func reserveCurrentNS() (func(), error) {
	ns := map[string]int{
		"ipc":  syscall.CLONE_NEWIPC,
		"mnt":  syscall.CLONE_NEWNS,
		"net":  syscall.CLONE_NEWNET,
		//TODO pid and user ns
		//"pid":  syscall.CLONE_NEWPID,
		//"user": syscall.CLONE_NEWUSER,
		"uts":  syscall.CLONE_NEWUTS,
	}
	fds := make(map[string]int)
	for name := range ns {
		fd, err := GetNSfd(fmt.Sprintf("/proc/self/ns/%s", name))
		if err != nil {
			return nil, fmt.Errorf("can't get %s fd: %v", name, err)
		}
		fds[name] = fd
	}
	return func() {
		for name, fd := range fds {
			if err := Setns(fd, ns[name]); err != nil {
				glog.Warningf("failed to setns %s, fd %d", name, fd)
			}
			if err := Close(fd); err != nil {
				glog.Warningf("failed to close ns %s, fd %d", name, fd)
			}
		}
	}, nil
}

// Close closes the fd
func Close(fd int) error {
	if err := syscall.Close(fd); err != nil {
		return err
	}
	return nil
}
