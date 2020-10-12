package libcontainer

const stdioFdCount = 3

type linuxContainer struct {
	id 		string
	root 		string
}

type State struct {
	BaseState

	CgroupPaths 	map[string]string 		`json:"cgroup_paths"`
}


type Container interface {
	BaseContainer

	NotifyOOM() (<-chan struct{}, error)
}


func (c *linuxContainer) Start(process *Process) error {
	// TODO process.Init

}

func (c *linuxContainer) start(process *Process) error {

}

func (c *linuxContainer) newParentProcess(p *Process) (parentProcess, error) {
	return nil, nil
}