package main

import (
	"flag"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/golang/glog"
)

var (
	rootDir = flag.String("root", "", "the root directory of container")
)

func main() {
	flag.Parse()
	defer glog.Flush()
	if *rootDir == "" {
		*rootDir = "."
	}
	tailArgs := flag.Args()
	if len(tailArgs) < 1 {
		glog.Fatal("command required")
	}
	glog.V(1).Infof("root %s, command %v", *rootDir, tailArgs)
	teardown := NewContainer()
	defer teardown()
	chback, err := chroot(*rootDir)
	if err != nil {
		glog.Fatal(err)
	}
	defer chback()
	if err := mountfs(); err != nil {
		glog.Fatal(err)
	}
	if err := syscall.Exec(tailArgs[0], tailArgs[0:], strings.Split("PATH=/sbin:/bin:/usr/sbin:/usr/bin", " ")); err != nil {
		glog.Fatal(err)
	}
}

func NewContainer() func() {
	runtime.LockOSThread()
	setback, err := reserveCurrentNS()
	if err != nil {
		glog.Fatal(err)
	}
	//TODO pid and user ns
	if err := syscall.Unshare(syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUTS); err != nil {
		glog.Fatal(err)
	}
	return func() {
		setback()
		runtime.UnlockOSThread()
	}
}

func chroot(root string) (func(), error) {
	if err := syscall.Mount("", "/", "", uintptr(syscall.MS_SLAVE|syscall.MS_REC), ""); err != nil {
		return nil, err
	}
	realRoot, err := os.Open("/")
	if err != nil {
		return nil, err
	}
	if err := syscall.Chdir(root); err != nil {
		return nil, err
	}
	if err := syscall.Chroot(root); err != nil {
		return nil, err
	}
	return func() {
		if err := syscall.Fchdir(int(realRoot.Fd())); err != nil {
			glog.Fatal(err)
		}
		if err := syscall.Chroot("."); err != nil {
			glog.Fatal(err)
		}
		realRoot.Close()
	}, nil
}

func mountfs() error {
	if err := os.MkdirAll("/proc", 0755); err != nil {
		glog.Error(err)
		return err
	}
	if err := os.MkdirAll("/sys", 0755); err != nil {
		glog.Error(err)
		return err
	}
	defaultMountFlags := uintptr(syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV)
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		glog.Error(err)
		return err
	}
	if err := syscall.Mount("sysfs", "/sys", "sysfs", uintptr(defaultMountFlags), ""); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
