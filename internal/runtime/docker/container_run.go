package docker

import (
	"context"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
)

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
