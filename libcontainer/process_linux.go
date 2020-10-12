package libcontainer

const createCgroupns = 0x80

type parentProcess interface {
	// pid returns the pid for the running process.
	pid() int

	// start starts the process exection
	start() error
}
