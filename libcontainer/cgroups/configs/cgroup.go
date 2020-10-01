package configs

type Cgroup struct {
	// Deprecated, use Path instead
	Name string `json:"name,omitempty"`

	// name of parent of cgroup or slice
	// Deprecated, use Path instead
	Parent string `json:"parent,omitempty`

	Path   string `json:"path"`

	// ScopePrefix describes prefix for the scope name
	ScopePrefix 	string 	`json:"scope_prefix"`

	// Paths represent the absolute cgroups paths to join
	// This takes precedence over Path
	Paths map[string]string

	// TODO Resources
}
