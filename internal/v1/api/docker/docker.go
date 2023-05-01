package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
)

type DockerRun struct {
	Name            *string
	Image           string
	Command         string
	PullImageAlways bool
}

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

	containerID, err = docker.containerCreateAndStart(ctx, "")
	if err != nil {
		return err
	}

	return docker.exec(ctx, containerID, autoStop)
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
func (docker *Docker) Run(ctx context.Context, run DockerRun) error {

	imageExists, err := docker.checkIfImageExists(ctx, run.Image)
	if err != nil {
		return err
	}

	if run.PullImageAlways {
		fmt.Println("'-pull-image-always' option is disabled!")
	}

	if !imageExists {
		err := docker.pullImage(ctx, run.Image)
		if err != nil {
			return err
		}
	}

	containerID, err := docker.containerCreateAndStart(ctx, "")
	if err != nil {
		return err
	}

	return docker.exec(ctx, containerID, true)
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

func (docker *Docker) checkIfImageExists(ctx context.Context, image string) (bool, error) {

	if len(strings.Split(image, ":")) != 2 {
		image = image + ":latest"
	}

	images, err := docker.cli.ImageList(ctx, types.ImageListOptions{
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

func (docker *Docker) exec(ctx context.Context, containerID string, autoStop bool) error {

	execID, err := docker.cli.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Privileged:   true,
		Tty:          true,
		Cmd:          []string{"/bin/bash"},
	})
	if err != nil {
		panic(err)
	}

	// Attach to the exec instance to read its output
	hijackedResp, err := docker.cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		panic(err)
	}
	defer hijackedResp.Close()

	// Start the exec instance
	err = docker.cli.ContainerExecStart(context.Background(), execID.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	// Copy input/output between the terminal and the container
	go func() {
		_, err = io.Copy(os.Stdout, hijackedResp.Reader)
		if err != nil && err != io.EOF {
			panic(err)
		}
	}()



	_, err = io.Copy(hijackedResp.Conn, os.Stdin)
	if err != nil && err != io.EOF {
		return err
	}

	if autoStop {
		return docker.Stop(ctx)
	}

	return nil
}

func (docker *Docker) pullImage(ctx context.Context, image string) error {

	out, err := docker.cli.ImagePull(ctx, image, types.ImagePullOptions{})
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

func (docker *Docker) containerCreateAndStart(ctx context.Context, image string) (string, error) {

	if image == "" {
		image = "envcontainer/envcontainer"
	}
	// Create the container
	containerResponse, err := docker.cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
		Tty: true,
		Cmd: []string{"/bin/bash"},
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
		return "", err
	}

	// Start the container
	err = docker.cli.ContainerStart(ctx, containerResponse.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return containerResponse.ID, nil
}
