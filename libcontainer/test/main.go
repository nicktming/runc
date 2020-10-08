package main

import (
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"fmt"
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

}

func main() {

}

