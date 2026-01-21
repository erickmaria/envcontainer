package types

import (
	"github.com/ErickMaria/envcontainer/internal/pkg/types"
)

type Envcontainer struct {
	Project   Project   `yaml:"project"`
	Container Container `yaml:"container"`
	// Labels       []string  `yaml:"labels"`
	AlwaysUpdate bool `yaml:"always_update"`
	AutoStop     bool `yaml:"auto_stop"`
	mountDir     string
	Mounts       []types.Mount `yaml:"mounts"`
}

func (envcontainer *Envcontainer) SetMountDir(mountDir string) {
	envcontainer.mountDir = mountDir
}

func (envcontainer *Envcontainer) GetMountDir() string {
	return envcontainer.mountDir
}
