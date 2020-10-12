package libcontainer

type Status int

const (
	Created Status = iota
	Running
	Pausing
	Paused
	Stopped
)

func (s Status) String() string {
	switch s {
	case Created:
		return "created"
	case Running:
		return "running"
	case Pausing:
		return "pausing"
	case Paused:
		return "paused"
	case Stopped:
		return "stopped"
	default:
		return "unknown"
	}
}

type BaseState struct {
	// ID is the container ID.
	ID 	string 		`json:"id"`
}

type BaseContainer interface {
	ID() 	string
}