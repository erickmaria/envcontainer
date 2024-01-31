package types

type ContainerOptions struct {
	User              string
	HomeDir           string
	AutoStop          bool
	ContainerName     string
	ImageName         string
	PullImageAlways   bool
	Commands          []string
	Ports             []string
	HostDirToBind     string
	Mounts            []string
	NoContainerSuffix bool
}
