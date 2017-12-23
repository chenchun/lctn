package lctn

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/chenchun/lctn/flags"
	"github.com/golang/glog"
)

func Parent() {
	args := []string{os.Args[0]}
	flag.Set("init", "true")
	flag.Visit(func(f *flag.Flag) {
		args = append(args, fmt.Sprintf("-%s", f.Name))
		//golang flag.go
		//Command line flag syntax:
		//-flag
		//-flag=x
		//-flag x  // non-boolean flags only
		if f.Value.String() != "true" && f.Value.String() != "false" {
			args = append(args, f.Value.String())
		}
	})
	args = append(args, flag.Args()...)
	glog.V(2).Infof("init args %v", args)
	cmd := exec.Cmd{
		Path: os.Args[0],
		Args: args,
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags: uintptr(syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER),
		},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Start(); err != nil {
		glog.Fatal(err)
	}
	go func() {
		if err := mapUidGid(cmd.Process.Pid); err != nil {
			glog.Fatal(err)
		}
	}()
	if err := cmd.Wait(); err != nil {
		glog.Fatal(err)
	}
}

func Child() {
	if err := Chroot(*flags.RootDir); err != nil {
		glog.Fatal(err)
	}
	if err := Mountfs(); err != nil {
		glog.Fatal(err)
	}
	tailArgs := flag.Args()
	glog.V(1).Infof("container command and args %v", tailArgs)
	if err := syscall.Exec(tailArgs[0], tailArgs, strings.Split("PATH=/sbin:/bin:/usr/sbin:/usr/bin", " ")); err != nil {
		glog.Fatal(err)
	}
}
