package libcontainer

import (
	"os"
	"fmt"
	"math"
	"io"
)

type processOperations interface {
	wait() (*os.ProcessState, error)
	signal(sig os.Signal) error
	pid() int
}

type Process struct {
	Args []string

	Env  []string

	User string

	Cwd  string

	ops processOperations
}

// Wait waits for the process to exit.
// Wait releases any resources associated with the Process
func (p Process) Wait() (*os.ProcessState, error) {
	if p.ops == nil {
		return nil, newGenericError(fmt.Errorf("invalid process"), NoProcessOps)
	}
	return p.ops.wait()
}

// Pid returns the process ID
func (p Process) Pid() (int, error) {
	// math.MinInt32 is returned here, because it's invalid value
	// for the kill() system call.
	if p.ops == nil {
		return math.MinInt32, newGenericError(fmt.Errorf("invalid process"), NoProcessOps)
	}
	return p.ops.pid(), nil
}

// Signal sends a signal to the Process.
func (p Process) Signal(sig os.Signal) error {
	if p.ops == nil {
		return newGenericError(fmt.Errorf("invalid process"), NoProcessOps)
	}
	return p.ops.signal(sig)
}

// IO holds the process's STDIO
type IO struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}
