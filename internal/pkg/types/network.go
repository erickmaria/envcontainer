package types

type Network struct {
	Name       string            `yaml:"name,omitempty"`
	External   bool              `yaml:"external,omitempty"`
	Driver     string            `yaml:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	IPAM       *IPAM             `yaml:"ipam,omitempty"`
	// Labels     map[string]string `yaml:"labels,omitempty"`
}

type IPAM struct {
	Driver string       `yaml:"driver,omitempty"`
	Config []IPAMConfig `yaml:"config,omitempty"`
}

type IPAMConfig struct {
	Subnet  string `yaml:"subnet,omitempty"`
	Gateway string `yaml:"gateway,omitempty"`
}
