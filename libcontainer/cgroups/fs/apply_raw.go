package fs

import (
	"sync"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"fmt"
	"os"
)

type cgroupData struct {
	root 		string
	innerPath 	string
	config		*configs
	pid 		int
}

var cgroupRootLock sync.Mutex
var cgroupRoot string

// Gets the cgroupRoot
func getCgroupRoot() (string, error) {
	cgroupRootLock.Lock()
	defer cgroupRootLock.Unlock()

	if cgroupRoot != "" {
		return cgroupRoot, nil
	}

	root, err := cgroups.FindCgroupMountpointDir()
	if err != nil {
		return "", err
	}
	fmt.Printf("root: %v\n", root)

	if _, err := os.Stat(root); err != nil {
		return "", err
	}

	cgroupRoot = root
	return cgroupRoot, nil
}