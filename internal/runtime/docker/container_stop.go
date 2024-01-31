package docker

import (
	"context"
	"fmt"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/docker/docker/api/types"
)

func (docker *Docker) Stop(ctx context.Context, options runtimeTypes.ContainerOptions) error {

	docker.addContainerSuffix(&options)

	for {
		container, err := docker.getContainer(ctx, options.ContainerName)

		if err != nil {
			return err
		}

		if container.ID == "" {
			fmt.Println("no containers found with name '" + options.ContainerName + "'")
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

		return nil
	}

}
