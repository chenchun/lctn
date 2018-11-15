package lctn

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/golang/glog"
)

func Chroot(root string) error {
	if err := syscall.Mount("", "/", "", uintptr(syscall.MS_PRIVATE|syscall.MS_REC), ""); err != nil {
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

func Mountfs(cgInfo *CgroupInfo) error {
	if err := os.MkdirAll("/proc", 0755); err != nil {
		return err
	}
	if err := os.MkdirAll("/sys", 0755); err != nil {
		return err
	}
	defaultMountFlags := uintptr(syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV)
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		return err
	}
	if err := syscall.Mount("sysfs", "/sys", "sysfs", uintptr(defaultMountFlags), ""); err != nil {
		return err
	}
	if err := syscall.Mount("cgroup", "/sys/fs/cgroup", "tmpfs", uintptr(defaultMountFlags), ""); err != nil {
		return err
	}
	for _, sub := range cgInfo.CgroupSubSystems {
		subsystemPath := filepath.Join(cgInfo.CgroupRoot, sub)
		if _, err := os.Stat(subsystemPath); err != nil {
			if err := os.MkdirAll(subsystemPath, 0755); err != nil {
				glog.Warning(err)
				continue
			}
		}
		if err := syscall.Mount("cgroup", subsystemPath, "cgroup", 0, sub); err != nil {
			glog.Warningf("failed to mount cgroup subsystem %s: %v", sub, err)
		}
	}
	return nil
}

func PrepareDevice(root string, bind bool) error {
	devPath := filepath.Join(root, "/dev")
	if err := os.MkdirAll(devPath, 0755); err != nil {
		return err
	}
	if !bind {
		if err := syscall.Mount("devtmpfs", devPath, "devtmpfs", 0, "rw,nosuid,relatime,size=6031164k,mode=755"); err != nil {
			return err
		}
	} else {
		nullPath := filepath.Join(root, "/dev/null")
		if _, err := os.Stat(nullPath); err != nil {
			if _, err := os.Create(nullPath); err != nil {
				glog.Error(err)
				return err
			}
		}
		if err := syscall.Mount("/dev/null", nullPath, "", syscall.MS_BIND, ""); err != nil && !os.IsNotExist(err) {
			glog.Error(err)
			return err
		}
	}
	return nil
}

func RemoveDevice(root string) {
	os.RemoveAll(filepath.Join(root, "/dev"))
}
