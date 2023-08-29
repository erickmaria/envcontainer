package docker

import (
	"context"
	"strings"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/docker/docker/client"
)

var shells = []string{"/bin/bash", "/bin/sh"}

type Docker struct {
	client *client.Client
}

func NewDocker() *Docker {

	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &Docker{
		client: client,
	}
}

func (docker *Docker) AlwaysUpdate(ctx context.Context, options runtimeTypes.BuildOptions) error {

	err := docker.Stop(ctx, options.ImageName)
	if err != nil {
		panic(err)
	}

	return docker.Build(ctx, options)

}

func (docker *Docker) addContainerSuffix(options *runtimeTypes.ContainerOptions) {

	pathSplit := strings.Split(options.HostDirToBind, "/")
	containerNameSuffix := pathSplit[len(pathSplit)-1]
	if containerNameSuffix != options.ContainerName {
		options.ContainerName = options.ContainerName + "-" + containerNameSuffix
	}
}
