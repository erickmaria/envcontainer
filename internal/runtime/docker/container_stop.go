package docker

import (
	"context"
	"fmt"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/docker/docker/api/types/container"
)

func (docker *Docker) Stop(ctx context.Context, options runtimeTypes.ContainerOptions) error {

	docker.addContainerSuffix(&options)

	for {
		getContainer, err := docker.getContainer(ctx, options.ContainerName)

		if err != nil {
			return err
		}

		if getContainer.ID == "" {
			fmt.Println("no containers found with name '" + options.ContainerName + "'")
			return nil
		}

		// Stopping the container
		fmt.Print("Stopping container ", getContainer.ID[:10], "... ")

		// Remove the container
		docker.client.ContainerRemove(ctx, getContainer.ID, container.RemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		})

		fmt.Println("Success!")

		return nil
	}

}
