package fs

import (
	"sync"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"fmt"
	"os"
	"github.com/opencontainers/runc/libcontainer/cgroups/configs"
	libcontainerUtils "github.com/opencontainers/runc/libcontainer/utils"
	"path/filepath"
	"encoding/json"
	"golang.org/x/sys/unix"
	"github.com/pkg/errors"
)

var (
	subsystems = subsystemSet{
		&CpuGroup{},
	}
)

var errSubsystemDoesNotExist = fmt.Errorf("cgroup: subsystem does not exist")

type subsystemSet []subsystem

func (s subsystemSet) Get(name string) (subsystem, error) {
	for _, ss := range s {
		if ss.Name() == name {
			return ss, nil
		}
	}
	return nil, errSubsystemDoesNotExist
}

type subsystem interface {
	Name() string
	//GetStats(path string, stats *cgroups.Stats) error
	Remove(*cgroupData) error
	Apply(*cgroupData) error
	Set(path string, cgroup *configs.Cgroup) error
}

type Manager struct {
	mu 		sync.Mutex
	Cgroups 	*configs.Cgroup
	Rootless 	bool
	Paths 		map[string]string
}

func isIgnorableError(rootless bool, err error) bool {
	// We do not ignore errors if we are root.
	if !rootless {
		return false
	}
	// Is it an ordinary EPERM?
	if os.IsPermission(errors.Cause(err)) {
		return true
	}

	// Try to handle other errnos.
	var errno error
	switch err := errors.Cause(err).(type) {
	case *os.PathError:
		errno = err.Err
	case *os.LinkError:
		errno = err.Err
	case *os.SyscallError:
		errno = err.Err
	}
	return errno == unix.EROFS || errno == unix.EPERM || errno == unix.EACCES
}

func (m *Manager) Apply(pid int) (err error) {
	if m.Cgroups == nil {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	var c = m.Cgroups
	d, err := getCgroupData(m.Cgroups, pid)
	if err != nil {
		return err
	}
	// TODO m.Paths & c.Paths

	for _,sys := range subsystems {
		p, err := d.path(sys.Name())
		if err != nil {
			if cgroups.IsNotFound(err) && sys.Name() != "devices" {
				continue
			}
			return err
 		}
		m.Paths[sys.Name()] = p
		if err := sys.Apply(d); err != nil {
			if isIgnorableError(m.Rootless, err) && m.Cgroups.Path == "" {
				delete(m.Paths, sys.Name())
				continue
			}
		}
		return err
	}
	return nil
}

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
	mnt, err :=  cgroups.FindCgroupMountpoint(raw.root, subsystem)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(raw.innerPath) {
		return filepath.Join(raw.root, filepath.Base(mnt), raw.innerPath), nil
	}

	parentPath, err := cgroups.GetOwnCgroupPath(subsystem)
	if err != nil {
		return "", err
	}
	pretty_json, _ := json.MarshalIndent(raw, "", "\t")
	fmt.Printf("=======>parentPath: %v, subsystem: %v, pretty_json: %v\n",
		parentPath, subsystem, string(pretty_json))
	return filepath.Join(parentPath, raw.innerPath), nil
}






















































