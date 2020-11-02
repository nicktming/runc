package main

import (
	"github.com/opencontainers/runc/libcontainer"
	"os"
	"log"
	"github.com/opencontainers/runc/libcontainer/configs"
	"syscall"
)

func main() {
	root, err := libcontainer.New("/var/lib/container", libcontainer.InitArgs(os.Args[0], "init"))
	if err != nil {
		log.Fatal(err)
	}
	config := &configs.Config{
		Rootfs: rootfs,
		Capabilities: []string{
			"CHOWN",
			"DAC_OVERRIDE",
			"FSETID",
			"FOWNER",
			"MKNOD",
			"NET_RAW",
			"SETGID",
			"SETUID",
			"SETFCAP",
			"SETPCAP",
			"NET_BIND_SERVICE",
			"SYS_CHROOT",
			"KILL",
			"AUDIT_WRITE",
		},
		Namespaces: configs.Namespaces([]configs.Namespace{
			{Type: configs.NEWNS},
			{Type: configs.NEWUTS},
			{Type: configs.NEWIPC},
			{Type: configs.NEWPID},
			{Type: configs.NEWNET},
		}),
		Cgroups: &configs.Cgroup{
			Name:            "test-container",
			Parent:          "system",
			AllowAllDevices: false,
			AllowedDevices:  configs.DefaultAllowedDevices,
		},

		Devices:  configs.DefaultAutoCreatedDevices,
		Hostname: "testing",
		Networks: []*configs.Network{
			{
				Type:    "loopback",
				Address: "127.0.0.1/0",
				Gateway: "localhost",
			},
		},
		Rlimits: []configs.Rlimit{
			{
				Type: syscall.RLIMIT_NOFILE,
				Hard: uint64(1024),
				Soft: uint64(1024),
			},
		},
	}

	container, err := root.Create("container-id", config)

	process := &libcontainer.Process{
		Args:   []string{"/bin/bash"},
		Env:    []string{"PATH=/bin"},
		User:   "daemon",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = container.Start(process)
	if err != nil {
		log.Fatal(err)
	}

	// wait for the process to finish.
	status, err := process.Wait()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status: %v\n", status)

	// destroy the container.
	container.Destroy()
}
