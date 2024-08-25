package docker

import (
	"context"
	"fmt"
	"sort"

	"github.com/ErickMaria/envcontainer/internal/runtime/types"
)

func (docker *Docker) List(ctx context.Context, options map[string]types.ContainerOptions) error {

	fmt.Printf("%-20s  %-20s %-20s\n", "CONTAINER NAME", "STATUS", "PATH")

	paths := make([]string, 0, len(options))
	for key := range options {
		paths = append(paths, key)
	}
	sort.Strings(paths)

	for _, path := range paths {
		getContainer, err := docker.getContainer(ctx, options[path].ContainerName)
		if err != nil {
			return err
		}
		fmt.Printf("%-20s %-20s  %-20s\n", options[path].ContainerName, getContainer.State, path)
	}

	return nil

}
