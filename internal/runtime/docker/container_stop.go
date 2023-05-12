package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
)

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
