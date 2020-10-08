package fs

import (
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"fmt"
)

type CpuGroup struct {

}

func (s *CpuGroup) Name() string {
	return "cpu"
}

func (s *CpuGroup) Apply(d *cgroupData) error {
	// We always want to join the cpu group, to allow fair cpu scheduling
	// on a container basis
	path, err := d.path("cpu")
	if err != nil && !cgroups.IsNotFound(err) {
		return err
	}
	fmt.Printf("=======>path: %v\n", path)
	return nil
}
