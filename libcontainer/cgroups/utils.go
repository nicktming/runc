package cgroups

import (
	"io"
	"bufio"
	"strings"
	"fmt"
	"os"
	"path/filepath"
)

const (
	CgroupNamePrefix = "name="
	CgroupProcesses  = "cgroup.procs"
)

func FindCgroupMountpoint(cgroupPath, subsystem string) (string, error) {
	mnt, _, err := FindCgroupMountpointAndRoot(cgroupPath, subsystem)
	return mnt, err
}

func FindCgroupMountpointAndRoot(cgroupPath, subsystem string) (string, string, error) {
	if !isSubsystemAvailable(subsystem) {
		return "", "", NewNotFoundError(subsystem)
	}
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	return findCgroupMountpointAndRootFromReader(f, cgroupPath, subsystem)
}

func findCgroupMountpointAndRootFromReader(reader io.Reader, cgroupPath, subsystem string) (string, string, error) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Fields(txt)
		if len(fields) < 5 {
			continue
		}
		if strings.HasPrefix(fields[4], cgroupPath) {
			for _, opt := range strings.Split(fields[len(fields)-1], ",") {
				if opt == subsystem {
					return fields[4], fields[3], nil
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", err
	}
	return "", "", NewNotFoundError(subsystem)
}

func isSubsystemAvailable(subsystem string) bool {
	cgroups, err := ParseCgroupFile("/proc/self/cgroup")
	if err != nil {
		return false
	}
	_, avail := cgroups[subsystem]
	return avail
}

func FindCgroupMountpointDir() (string, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		fields := strings.Split(text, " ")

		index := strings.Index(text, " - ")
		fmt.Printf("text[index + 3L]=%v\n", text[index+3:])
		postSeparatorFields := strings.Fields(text[index+3:])
		numPostFields := len(postSeparatorFields)

		if numPostFields == 0 {
			return "", fmt.Errorf("Found no fields post '-' in %q", text)
		}

		if postSeparatorFields[0] == "cgroup" {
			// Check that the mount is properly formatted.
			if numPostFields < 3 {
				return "", fmt.Errorf("Error found less than 3 fields post '-' in %q", text)
			}

			return filepath.Dir(fields[4]), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", NewNotFoundError("cgroup")
}

func ParseCgroupFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseCgroupFromReader(f)
}

func parseCgroupFromReader(r io.Reader) (map[string]string, error) {
	s := bufio.NewScanner(r)
	cgroups := make(map[string]string)

	for s.Scan() {
		text := s.Text()

		parts := strings.SplitN(text, ":", 3)
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid cgroup entry: must contain at least two colons: %v", text)
		}

		for _, subs := range strings.Split(parts[1], ",") {
			cgroups[subs] = parts[2]
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return cgroups, nil
}

func GetOwnCgroup(subsystem string) (string, error) {
	cgroups, err := ParseCgroupFile("/proc/self/cgroup")
	if err != nil {
		return "", err
	}
	return getControllerPath(subsystem, cgroups)
}

func GetOwnCgroupPath(subsystem string) (string, error) {
	cgroup, err := GetOwnCgroup(subsystem)
	if err != nil {
		return "", err
	}
	return getCgroupPathHelper(subsystem, cgroup)
}

func getCgroupPathHelper(subsystem, cgroup string) (string, error) {
	mnt, root, err := FindCgroupMountpointAndRoot("", subsystem)
	if err != nil {
		return "", err
	}

	// This is needed for nested containers, because in /proc/self/cgroup we
	// see paths from host, which don't exist in container.
	relCgroup, err := filepath.Rel(root, cgroup)
	if err != nil {
		return "", err
	}

	return filepath.Join(mnt, relCgroup), nil
}

func getControllerPath(subsystem string, cgroups map[string]string) (string, error) {

	if p, ok := cgroups[subsystem]; ok {
		return p, nil
	}

	if p, ok := cgroups[CgroupNamePrefix+subsystem]; ok {
		return p, nil
	}

	return "", NewNotFoundError(subsystem)
}






































