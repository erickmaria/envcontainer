package docker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

func (docker *Docker) Start(ctx context.Context, options runtimeTypes.ContainerOptions) error {

	if options.ImageName == "" {
		options.ImageName = "envcontainer/" + options.ContainerName

		imageExists, err := docker.checkIfImageExists(ctx, options.ImageName)
		if err != nil {
			return err
		}
		if !imageExists {
			return errors.New("no such image try run 'build' command")
		}

	}

	docker.addContainerSuffix(&options)

	container, err := docker.getContainer(ctx, options.ContainerName)
	if err != nil {
		return err
	}

	if container.ID != "" {
		options.Commands = []string{container.Command}
		return docker.exec(ctx, container.ID, options)
	}

	return docker.tryCreateAndStartContainer(ctx, options)
}

func (docker *Docker) containerCreateAndStart(ctx context.Context, options runtimeTypes.ContainerOptions) error {

	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	if len(options.Ports) != 0 {
		for _, port := range options.Ports {

			splitPort := strings.Split(strings.Trim(port, " "), ":")

			bindPort := nat.PortBinding{}
			bindPort.HostIP = "0.0.0.0"

			if len(splitPort) == 2 {
				exposedPorts[nat.Port(splitPort[0])] = struct{}{}
				bindPort.HostPort = splitPort[1]
				portBindings[nat.Port(splitPort[0])] = []nat.PortBinding{bindPort}
				continue
			}

			exposedPorts[nat.Port(port)] = struct{}{}
			bindPort.HostPort = port
			portBindings[nat.Port(port)] = []nat.PortBinding{bindPort}

		}
	}

	mounts := options.Mounts
	mounts = append(mounts, options.HostDirToBind+":/home/"+options.ContainerName)

	// Create the container
	containerResponse, err := docker.client.ContainerCreate(ctx, &container.Config{
		User:         options.User,
		WorkingDir:   "/home/" + options.ContainerName,
		Image:        options.ImageName,
		ExposedPorts: exposedPorts,
		Tty:          true,
		Cmd:          options.Commands,
	}, &container.HostConfig{
		PortBindings: portBindings,
		Binds:        mounts,
	}, &network.NetworkingConfig{}, nil, options.ContainerName)
	if err != nil {
		return err
	}

	// Start the container
	err = docker.client.ContainerStart(ctx, containerResponse.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Print("Error to start container, ")
		docker.Stop(ctx, options.ContainerName)
		return err
	}

	return docker.exec(ctx, containerResponse.ID, options)
}

func (docker *Docker) tryCreateAndStartContainer(ctx context.Context, options runtimeTypes.ContainerOptions) error {

	if len(options.Commands) > 0 && options.Commands[0] != "" {
		return docker.containerCreateAndStart(ctx, options)
	}
	var err error
	for _, shell := range shells {
		options.Commands = []string{shell}
		if err = docker.containerCreateAndStart(ctx, options); err == nil {
			return nil
		}
	}

	return err
}

func (docker *Docker) getContainer(ctx context.Context, containerName string) (types.Container, error) {

	containers, err := docker.client.ContainerList(ctx, types.ContainerListOptions{
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
