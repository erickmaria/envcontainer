package types

type Mount struct {
	Source   string `yaml:"source"`
	Target   string `yaml:"target"`
	Type     string `yaml:"type"`
	Readonly bool   `yaml:"readonly"`
	// Labels   map[string]string `yaml:"labels,omitempty"`
}
