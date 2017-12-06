package lctn

import (
	"os"
	"syscall"
)

func Chroot(root string) error {
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

func Mountfs() error {
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
	return nil
}
