package fs

import (
	"sync"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"fmt"
	"os"
	"github.com/opencontainers/runc/libcontainer/cgroups/configs"
	libcontainerUtils "github.com/opencontainers/runc/libcontainer/utils"
	"github.com/opencontainers/runtime-spec/specs-go"
	"path/filepath"
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

func getCgroupData(c *configs.Cgroup, pid int) (*cgroupData, error) {
	root, err := getCgroupRoot()
	if err != nil {
		return nil, err
	}

	if (c.Name != "" || c.Parent != "") && c.Path != "" {
		return nil, fmt.Errorf("cgroup: either Path or Name and Parent should be used")
	}

	cgPath := libcontainerUtils.CleanPath(c.Path)
	cgParent := libcontainerUtils.CleanPath(c.Parent)
	cgName := libcontainerUtils.CleanPath(c.Name)

	innerPath := cgPath
	if innerPath == "" {
		innerPath = filepath.Join(cgParent, cgName)
	}

	return &cgroupData{
		root: 		root,
		innerPath: 	innerPath,
		config:		c,
		pid: 		pid,
	}, nil
}

func (raw *cgroupData) path(subsystem string) (string, error) {

}






















































