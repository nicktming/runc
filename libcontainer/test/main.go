package main

import (
	"fmt"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/configs"
)

func test1() {
	cgs, err := cgroups.ParseCgroupFile("/proc/self/cgroup")
	if err != nil {
		fmt.Errorf("err: %v", err)
		return
	}
	for k, v := range cgs {
		fmt.Printf("%v:%v\n", k, v)
	}

	root, err := cgroups.FindCgroupMountpointDir()
	if err != nil {
		fmt.Printf("FindCgroupMountpointDir with error: %v\n", root)
		return
	}
	fmt.Printf("root: %v\n", root)
}

func test2() {
	m := fs.Manager {
		Cgroups: &configs.Cgroup{
			Paths: 		make(map[string]string, 10),
		},
	}
	m.Apply(1)
}

func main() {
	test2()
}

