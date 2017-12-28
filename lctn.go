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
	reader, writer, err := os.Pipe()
	if err != nil {
		glog.Fatalf("can't create a pipe: %v", err)
	}
	defer writer.Close()
	cmd := exec.Cmd{
		Path: os.Args[0],
		Args: args,
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags: uintptr(syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER),
		},
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		ExtraFiles: []*os.File{reader},
	}
	if err := cmd.Start(); err != nil {
		glog.Fatal(err)
	}
	go func() {
		if err := mapUidGid(cmd.Process.Pid); err != nil {
			glog.Fatal(err)
		}
		// write 'd' (short for done) message to child process
		if _, err := writer.Write([]byte{'d'}); err != nil {
			glog.Fatalf("failed to write into pipe: %v", err)
		}
	}()
	if err := cmd.Wait(); err != nil {
		glog.Fatal(err)
	}
}

func Child() {
	root := *flags.RootDir
	if err := PrepareDevice(root); err != nil {
		glog.Fatal(err)
	}
	if err := Chroot(root); err != nil {
		glog.Fatal(err)
	}
	if err := Mountfs(); err != nil {
		glog.Fatal(err)
	}
	reader := os.NewFile(3, "parent-pipe")

	msg := make([]byte, 1)
	if _, err := reader.Read(msg); err != nil {
		glog.Fatalf("failed to read from parent: %v", err)
	}
	glog.V(2).Infof("read from parent: %s", string(msg))
	if err := reader.Close(); err != nil {
		glog.Warningf("failed to close read pipe: %v", err)
	}
	tailArgs := flag.Args()
	glog.V(1).Infof("container command and args %v", tailArgs)
	if err := syscall.Exec(tailArgs[0], tailArgs, strings.Split("PATH=/sbin:/bin:/usr/sbin:/usr/bin", " ")); err != nil {
		glog.Fatal(err)
	}
}
