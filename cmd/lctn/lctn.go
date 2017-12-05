package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/chenchun/lctn/flags"
	"github.com/golang/glog"
)

func main() {
	flags.InitFlags()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		glog.Fatal(err)
	}
	cmd := exec.Cmd{
		Path: filepath.Join(dir, "init"),
		Args: os.Args[1:],
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags: uintptr(syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID),
		},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Start(); err != nil {
		glog.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		glog.Fatal(err)
	}
}
