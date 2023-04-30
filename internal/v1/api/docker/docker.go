package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
)

type Docker struct {
	cli *client.Client
}

func NewDocker() *Docker {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &Docker{
		cli: cli,
	}
}

func (docker *Docker) Build(ctx context.Context) error {

	buildCtx, err := archive.TarWithOptions(".envcontainer", &archive.TarOptions{})
	if err != nil {
		return err
	}

	imageBuildResponse, err := docker.cli.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Tags: []string{"envcontainer/envcontainer"},
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

func (docker *Docker) Start(ctx context.Context, autoStop bool) error {

	containerID, err := docker.getContainerID(ctx, "envcontainer")
	if err != nil {
		return err
	}

	if containerID != "" {
		return docker.exec(ctx, containerID, autoStop)
	}

	// Create the container
	containerResponse, err := docker.cli.ContainerCreate(ctx, &container.Config{
		Image: "envcontainer/envcontainer",
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
		Tty: true,
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "8080",
				},
			},
		},
	}, &network.NetworkingConfig{}, nil, "envcontainer")
	if err != nil {
		return err
	}

	// Start the container
	err = docker.cli.ContainerStart(ctx, containerResponse.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return docker.exec(ctx, containerResponse.ID, autoStop)
}

func (docker *Docker) Stop(ctx context.Context) error {

	containerID, err := docker.getContainerID(ctx, "envcontainer")
	if err != nil {
		return err
	}

	if containerID == "" {
		fmt.Println("Not container running")
		return nil
	}

	// Stopping the container
	fmt.Print("Stopping container ", containerID[:10], "... ")

	// Remove the container
	docker.cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})

	fmt.Println("Success!")
	time.Sleep(1 * time.Second)

	return nil
}

func (docker *Docker) getContainerID(ctx context.Context, containerName string) (string, error) {

	containers, err := docker.cli.ContainerList(ctx, types.ContainerListOptions{
		Limit: 1,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: "envcontainer",
		}),
	})

	if len(containers) == 0 {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return containers[0].ID, nil
}

func (docker *Docker) exec(ctx context.Context, containerID string, autoStop bool) error {

	execID, err := docker.cli.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Privileged:   true,
		Tty:          true,
		Cmd: []string{
			"/bin/bash",
		},
	})
	if err != nil {
		panic(err)
	}

	// Attach to the exec instance to read its output
	execResp, err := docker.cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		panic(err)
	}
	defer execResp.Close()

	// Copy input/output between the terminal and the container
	go func() {
		_, err = io.Copy(os.Stdout, execResp.Reader)
		if err != nil && err != io.EOF {
			panic(err)
		}
	}()

	_, err = io.Copy(execResp.Conn, os.Stdin)
	if err != nil && err != io.EOF {
		panic(err)
	}

	if autoStop {
		return docker.Stop(ctx)
	}

	return nil
}
