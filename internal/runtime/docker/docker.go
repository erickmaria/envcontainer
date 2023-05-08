package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
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

func (docker *Docker) Build(ctx context.Context, options runtimeTypes.BuildOptions) error {

	buildCtx, err := archive.TarWithOptions("./", &archive.TarOptions{})
	if err != nil {
		return err
	}

	imageBuildResponse, err := docker.client.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Tags:       []string{"envcontainer/" + options.ImageName},
		Dockerfile: options.Dockerfile,
	})
	if err != nil {
		return err
	}
	defer imageBuildResponse.Body.Close()

	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return err
	}

	return nil
}

func (docker *Docker) Start(ctx context.Context, options runtimeTypes.ContainerOptions) error {

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

func (docker *Docker) Stop(ctx context.Context, containerName string) error {

	container, err := docker.getContainer(ctx, containerName)
	if err != nil {
		return err
	}

	if container.ID == "" {
		fmt.Println("Not container running")
		return nil
	}

	// Stopping the container
	fmt.Print("Stopping container ", container.ID[:10], "... ")

	// Remove the container
	docker.client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})

	fmt.Println("Success!")
	time.Sleep(1 * time.Second)

	return nil
}

func (docker *Docker) Run(ctx context.Context, options runtimeTypes.ContainerOptions) error {

	imageExists, err := docker.checkIfImageExists(ctx, options.ImageName)
	if err != nil {
		return err
	}

	if !imageExists {
		err := docker.pullImage(ctx, options.ImageName)
		if err != nil {
			return err
		}
	}

	return docker.tryCreateAndStartContainer(ctx, options)
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

func (docker *Docker) checkIfImageExists(ctx context.Context, image string) (bool, error) {

	if len(strings.Split(image, ":")) != 2 {
		image = image + ":latest"
	}

	images, err := docker.client.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: image,
		}),
	})
	if err != nil {
		return false, err
	}

	if len(images) == 0 {
		return false, nil
	}

	return true, nil
}

func (docker *Docker) exec(ctx context.Context, containerID string, options runtimeTypes.ContainerOptions) error {

	execID, err := docker.client.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Privileged:   true,
		Tty:          true,
		Cmd:          options.Commands,
	})
	if err != nil {
		return err
	}

	// Attach to the exec instance to read its output
	execResp, err := docker.client.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return err
	}
	defer execResp.Close()

	// Start the exec instance
	err = docker.client.ContainerExecStart(context.Background(), execID.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	cmd := exec.Command("docker", "exec", "-it", containerID, options.Commands[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// Copy input/output between the terminal and the container
	// go func() {
	// 	_, err = io.Copy(os.Stdout, execResp.Reader)
	// 	if err != nil && err != io.EOF {
	// 		panic(err)
	// 	}
	// }()

	// _, err = io.Copy(execResp.Conn, os.Stdin)
	// if err != nil && err != io.EOF {
	// 	return err
	// }

	if options.AutoStop {
		return docker.Stop(ctx, options.ContainerName)
	}

	return nil
}

func (docker *Docker) pullImage(ctx context.Context, image string) error {

	out, err := docker.client.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return err
	}

	return nil
}

func (docker *Docker) containerCreateAndStart(ctx context.Context, options runtimeTypes.ContainerOptions) error {

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
	// Create the container
	containerResponse, err := docker.client.ContainerCreate(ctx, &container.Config{
		User:       options.User,
		WorkingDir: "/home/" + options.ContainerName,
		Image:      options.ImageName,
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
		Tty: true,
		Cmd: options.Commands,
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "8080",
				},
			},
		},
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
