package lctn

import (
	"os"
	"path/filepath"
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

func PrepareDevice(root string) error {
	if err := os.MkdirAll("/dev", 0755); err != nil {
		return err
	}
	null := filepath.Join(root, "/dev/null")
	if _, err := os.Create(null); err != nil {
		return err
	}
	if err := syscall.Mount("/dev/null", null, "", syscall.MS_BIND, ""); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
