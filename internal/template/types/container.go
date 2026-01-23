package types

import "github.com/ErickMaria/envcontainer/internal/pkg/types"

type Container struct {
	Shell       string          `yaml:"shell"`
	Ports       []string        `yaml:"ports"`
	Build       string          `yaml:"build"`
	NetworkMode string          `yaml:"network_mode"`
	Networks    []types.Network `yaml:"networks"`
}
