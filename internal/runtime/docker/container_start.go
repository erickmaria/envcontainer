package docker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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

		inspect, err := docker.client.ContainerInspect(ctx, container.ID)
		if err != nil {
			return err
		}
		
		if (inspect.State.Status == "paused" || inspect.State.Status == "exited"){
			docker.tryStart(ctx, runtimeTypes.ContainerStartInfo{
				Name: inspect.Name,
				ID: inspect.ID,
			}, types.ContainerStartOptions{})
		}

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
	docker.tryStart(ctx, runtimeTypes.ContainerStartInfo{
		Name: options.ContainerName,
		ID: containerResponse.ID,
	}, types.ContainerStartOptions{})

	return docker.exec(ctx, containerResponse.ID, options)
}

func (docker *Docker) tryStart(ctx context.Context, info runtimeTypes.ContainerStartInfo, options types.ContainerStartOptions) error {
	err := docker.client.ContainerStart(ctx, info.ID, options)
	if err != nil {
		fmt.Print("Error to start container, ")
		docker.Stop(ctx, runtimeTypes.ContainerOptions{
			ContainerName: info.Name,
		})
		return err
	}
	return nil
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
