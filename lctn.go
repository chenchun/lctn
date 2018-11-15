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

const CLONE_NEWCGROUP = 0x02000000

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
	if err := os.Chown(*flags.RootDir, os.Getuid(), os.Getgid()); err != nil {
		glog.Fatalf("failed to chown rootfs: %v", err)
	}
	reader, writer, err := os.Pipe()
	if err != nil {
		glog.Fatalf("can't create a pipe: %v", err)
	}
	defer writer.Close()
	cmd := exec.Cmd{
		Path: os.Args[0],
		Args: args,
		SysProcAttr: &syscall.SysProcAttr{
			Setsid: true,
			Cloneflags: uintptr(syscall.CLONE_NEWNS |
				syscall.CLONE_NEWNET |
				syscall.CLONE_NEWIPC |
				syscall.CLONE_NEWUTS |
				syscall.CLONE_NEWPID |
				syscall.CLONE_NEWUSER |
				CLONE_NEWCGROUP),
			UidMappings: []syscall.SysProcIDMap{{ContainerID: 0, HostID: *flags.Uid, Size: 1}},
			GidMappings: []syscall.SysProcIDMap{{ContainerID: 0, HostID: *flags.Gid, Size: 1}},
		},
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		ExtraFiles: []*os.File{reader},
	}
	if err := cmd.Start(); err != nil {
		glog.Fatal(err)
	}
	cgPath := *flags.CgroupPath
	// echo cgroup pid must be done by parent process
	cgInfo, err := FindCgroupInfo()
	if err != nil {
		glog.Fatalf("failed to find cgroup info: %v", err)
	}
	if err := EnsureCgroup(cgPath, cgInfo); err != nil {
		glog.Fatalf("failed to ensure cgroup: %v", err)
	}
	defer RemoveCgroup(cgPath, cgInfo)
	defer RemoveDevice(*flags.RootDir)
	// tell child to mount cgroup now
	if _, err := writer.Write([]byte("d")); err != nil {
		glog.Fatalf("failed to write into pipe: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		glog.Fatal(err)
	}
}

func Child() {
	glog.V(3).Infof("uid: %d, gid: %d", os.Getuid(), os.Getgid())
	data, err := exec.Command("/sbin/getpcaps", fmt.Sprintf("%d", os.Getpid())).CombinedOutput()
	glog.V(3).Infof("getpcats: %s, err: %v", string(data), err)
	msg := make([]byte, 1)
	reader := os.NewFile(3, "parent-pipe")
	cgInfo, err := FindCgroupInfo()
	if err != nil {
		glog.Fatalf("failed to find cgroup info: %v", err)
	}
	root := *flags.RootDir
	// TODO pass whether new user namespace from parent
	if err := PrepareDevice(root, true); err != nil {
		glog.Fatal(err)
	}
	if err := Chroot(root); err != nil {
		glog.Fatal(err)
	}
	// wait for mount cgroup
	wait(reader, msg)
	if err := Mountfs(cgInfo); err != nil {
		glog.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		glog.Warningf("failed to close read pipe: %v", err)
	}
	tailArgs := flag.Args()
	glog.V(1).Infof("container command and args %v", tailArgs)
	if err := syscall.Exec(tailArgs[0], tailArgs, strings.Split("PATH=/sbin:/bin:/usr/sbin:/usr/bin", " ")); err != nil {
		glog.Fatal(err)
	}
}

func wait(reader *os.File, msg []byte) error {
	if _, err := reader.Read(msg); err != nil {
		return fmt.Errorf("failed to read from parent: %v", err)
	}
	return nil
}
