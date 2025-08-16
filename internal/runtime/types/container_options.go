package types

import "github.com/ErickMaria/envcontainer/internal/pkg/types"

type ContainerOptions struct {
	Shell             string
	HomeDir           string
	AutoStop          bool
	ContainerName     string
	ImageName         string
	PullImageAlways   bool
	Commands          []string
	Ports             []string
	HostDirToBind     string
	DefaultMountDir   string
	Mounts            []types.Mount
	Networks          []types.Network
	NoContainerSuffix bool
}
