package docker

import (
	"context"
	"strings"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
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

	err := docker.Stop(ctx, runtimeTypes.ContainerOptions{
		ContainerName: options.ImageName,
	})
	if err != nil {
		panic(err)
	}

	return docker.Build(ctx, options)

}

func (docker *Docker) addContainerSuffix(options *runtimeTypes.ContainerOptions) {

	if !options.NoContainerSuffix {
		pathSplit := strings.Split(options.HostDirToBind, "/")
		containerNameSuffix := pathSplit[len(pathSplit)-1]
		options.ContainerName = strings.ToLower(options.ContainerName + "-" + containerNameSuffix)
	}
}

func (docker *Docker) getContainer(ctx context.Context, containerName string) (types.Container, error) {

	containers, err := docker.client.ContainerList(ctx, container.ListOptions{
		Limit: 1,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: containerName,
		}),
	})

	if len(containers) == 0 {
		return types.Container{}, nil
	}

	if err != nil {
		return types.Container{}, err
	}

	return containers[0], nil
}
