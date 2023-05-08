package types

type ContainerOptions struct {
	AutoStop        bool
	ContainerName   string
	ImageName       string
	PullImageAlways bool
	Commands        []string
	Ports           []string
}
