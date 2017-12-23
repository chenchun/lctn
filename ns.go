package lctn

import (
	"fmt"
	"io/ioutil"
)

func mapUidGid(pid int) error {
	if err := ioutil.WriteFile(fmt.Sprintf("/proc/%d/uid_map", pid), []byte("0 1000 100"), 0600); err != nil {
		return fmt.Errorf("failed to write into uid_map: %v", err)
	}
	if err := ioutil.WriteFile(fmt.Sprintf("/proc/%d/gid_map", pid), []byte("0 1000 100"), 0600); err != nil {
		return fmt.Errorf("failed to write into gid_map: %v", err)
	}
	return nil
}
