package main

import (
	"flag"
	"os"
	"strings"
	"syscall"

	"github.com/chenchun/lctn/flags"
	"github.com/golang/glog"
)

func main() {
	flags.InitFlags()
	if err := chroot(*flags.RootDir); err != nil {
		glog.Fatal(err)
	}
	if err := mountfs(); err != nil {
		glog.Fatal(err)
	}
	tailArgs := flag.Args()
	if err := syscall.Exec(tailArgs[0], tailArgs[0:], strings.Split("PATH=/sbin:/bin:/usr/sbin:/usr/bin", " ")); err != nil {
		glog.Fatal(err)
	}
}

func chroot(root string) error {
	if err := syscall.Mount("", "/", "", uintptr(syscall.MS_SLAVE|syscall.MS_REC), ""); err != nil {
		return err
	}
	if err := syscall.Chdir(root); err != nil {
		return err
	}
	if err := syscall.Chroot(root); err != nil {
		return err
	}
	return nil
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
