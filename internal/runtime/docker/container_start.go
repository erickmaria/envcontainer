package docker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

func (docker *Docker) Up(ctx context.Context, options runtimeTypes.ContainerOptions, code bool, port string) error {

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

	getContainer, err := docker.getContainer(ctx, options.Labels)
	if err != nil {
		return err
	}

	if getContainer.ID != "" {

		inspect, err := docker.client.ContainerInspect(ctx, getContainer.ID)
		if err != nil {
			return err
		}

		if inspect.State.Status == "paused" || inspect.State.Status == "exited" {
			docker.tryStart(ctx, runtimeTypes.ContainerStartInfo{
				Name: inspect.Name,
				ID:   inspect.ID,
			}, container.StartOptions{})
		}

		options.Commands = []string{getContainer.Command}

		if code {
			return docker.code(ctx, getContainer.ID, port, options)
		}

		return docker.exec(ctx, getContainer.ID, options)
	}

	return docker.tryCreateAndStartContainer(ctx, options, code, port)
}

func (docker *Docker) containerCreateAndStart(ctx context.Context, options runtimeTypes.ContainerOptions, code bool, port string) error {

	var err error
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

	networkIDs := []string{}
	if len(options.Networks) > 0 {

		networkIDs, err = docker.createNetwork(ctx, options.Networks, options.Labels)
		if err != nil {
			return err
		}
	}

	mounts := docker.buildMount(options.DefaultMountDir, options.Mounts, options.Labels)

	bindProject := options.HostDirToBind + ":/home/" + options.ContainerName

	// Create the container
	containerResponse, err := docker.client.ContainerCreate(ctx, &container.Config{
		WorkingDir:   "/home/" + options.ContainerName,
		Image:        options.ImageName,
		ExposedPorts: exposedPorts,
		Tty:          true,
		Hostname:     options.ContainerName,
		Labels:       options.Labels,
	}, &container.HostConfig{
		PortBindings: portBindings,
		Binds:        []string{bindProject},
		Mounts:       mounts,
		NetworkMode:  container.NetworkMode(options.NetworkMode),
	}, &network.NetworkingConfig{}, nil, options.ContainerName)
	if err != nil {
		return err
	}

	// Start the container
	docker.tryStart(ctx, runtimeTypes.ContainerStartInfo{
		Name: options.ContainerName,
		ID:   containerResponse.ID,
	}, container.StartOptions{})

	if len(networkIDs) > 0 {
		for _, netId := range networkIDs {
			err := docker.client.NetworkConnect(ctx, netId, containerResponse.ID, &network.EndpointSettings{})
			if err != nil {
				fmt.Println("Error to connect network:", netId[:12], err)
			}
		}
	}

	if code {
		return docker.code(ctx, containerResponse.ID, port, options)
	}

	return docker.exec(ctx, containerResponse.ID, options)
}

func (docker *Docker) tryStart(ctx context.Context, info runtimeTypes.ContainerStartInfo, options container.StartOptions) error {
	err := docker.client.ContainerStart(ctx, info.ID, options)
	if err != nil {
		fmt.Print("Error to start container, ")
		docker.Down(ctx, runtimeTypes.ContainerOptions{
			ContainerName: info.Name,
		})
		return err
	}

	return nil
}

func (docker *Docker) tryCreateAndStartContainer(ctx context.Context, options runtimeTypes.ContainerOptions, code bool, port string) error {

	if len(options.Commands) > 0 && options.Commands[0] != "" {
		return docker.containerCreateAndStart(ctx, options, code, port)
	}
	var err error
	for _, shell := range shells {
		options.Commands = []string{shell}
		if err = docker.containerCreateAndStart(ctx, options, code, port); err == nil {
			return err
		}
	}

	return err
}
