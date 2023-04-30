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

func (docker *Docker) Run(ctx context.Context) error {

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

	return nil
}

func (docker *Docker) Stop(ctx context.Context) error {

	containers, err := docker.cli.ContainerList(ctx, types.ContainerListOptions{
		Limit: 1,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: "envcontainer",
		}),
	})
	if err != nil {
		return err
	}
	if len(containers) == 0 {
        panic("Container not found")
    }

	containerID := containers[0].ID

	// Stopping the container
	fmt.Print("Stopping container ", containerID[:10], "... ")
	
	// err = docker.cli.ContainerStop(ctx, containerID, container.StopOptions{})
	// if err != nil {
	// 	return err
	// }

	// Remove the container
	docker.cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force: true,
		RemoveVolumes: true,
	})

	fmt.Println("Success!")
	time.Sleep(1 * time.Second)

	return nil
}
