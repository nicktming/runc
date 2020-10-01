package main

import (
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"fmt"
)

func main() {
	cgroups, err := cgroups.ParseCgroupFile("/proc/self/cgroup")
	if err != nil {
		fmt.Errorf("err: %v", err)
		return
	}
	for k, v := range cgroups {
		fmt.Printf("%v:%v\n", k, v)
	}
}

