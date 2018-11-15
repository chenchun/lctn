package lctn

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type CgroupInfo struct {
	CgroupRoot       string
	CgroupSubSystems []string
}

func FindCgroupInfo() (*CgroupInfo, error) {
	mounts, err := ioutil.ReadFile("/proc/self/mounts")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc/self/mountinfo: %v", err)
	}
	var (
		device, mountpoint, fileSystemType string
		cgInfo                             CgroupInfo
	)
	sc := bufio.NewScanner(bytes.NewReader(mounts))
	for sc.Scan() {
		line := sc.Text()
		if n, err := fmt.Sscanf(line, "%s %s %s", &device, &mountpoint, &fileSystemType); n == 3 && err == nil {
			if fileSystemType == "cgroup" {
				if path.Base(mountpoint) == "systemd" {
					continue
				}
				cgInfo.CgroupRoot = path.Dir(mountpoint)
				cgInfo.CgroupSubSystems = append(cgInfo.CgroupSubSystems, path.Base(mountpoint))
			}
		}
	}
	if cgInfo.CgroupRoot == "" {
		return nil, fmt.Errorf("failed to find cgroup root path")
	}
	return &cgInfo, nil
}

func EnsureCgroup(path string, cgInfo *CgroupInfo) error {
	pidBytes := []byte(strconv.Itoa(os.Getpid()))
	if path != "" {
		for _, sub := range cgInfo.CgroupSubSystems {
			if sub == "cpuset" {
				dir := filepath.Join(cgInfo.CgroupRoot, sub, path)
				if err := os.MkdirAll(dir, 0644); err != nil {
					glog.Warning(err)
					break
				}
				cpusetCpus, err := ioutil.ReadFile(filepath.Join(cgInfo.CgroupRoot, "cpuset", "cpuset.cpus"))
				if err != nil {
					glog.Warning(err)
					break
				}
				cpusetMems, err := ioutil.ReadFile(filepath.Join(cgInfo.CgroupRoot, "cpuset", "cpuset.mems"))
				if err != nil {
					glog.Warning(err)
					break
				}
				if err := ioutil.WriteFile(filepath.Join(dir, "cpuset.cpus"), cpusetCpus, 0); err != nil {
					glog.Warning(err)
				}
				if err := ioutil.WriteFile(filepath.Join(dir, "cpuset.mems"), cpusetMems, 0); err != nil {
					glog.Warning(err)
				}
				break
			}
		}
	}
	for _, sub := range cgInfo.CgroupSubSystems {
		dir := filepath.Join(cgInfo.CgroupRoot, sub, path)
		if err := os.MkdirAll(dir, 0600); err != nil {
			continue
		}
		if err := ioutil.WriteFile(filepath.Join(cgInfo.CgroupRoot, sub, path, "cgroup.procs"), pidBytes, 0); err != nil {
			glog.Warning(err)
		}
	}
	return nil
}

func RemoveCgroup(path string, cgInfo *CgroupInfo) {
	if path != "" {
		pidBytes := []byte(strconv.Itoa(os.Getpid()))
		for _, sub := range cgInfo.CgroupSubSystems {
			// move pid to root, or we'll get "device or resource busy" on removeDir
			if err := ioutil.WriteFile(filepath.Join(cgInfo.CgroupRoot, sub, "cgroup.procs"), pidBytes, 0); err != nil {
				glog.Warning(err)
			}
			if err := os.RemoveAll(filepath.Join(cgInfo.CgroupRoot, sub, path)); err != nil {
				glog.Warning(err)
			}
		}
	}
}
