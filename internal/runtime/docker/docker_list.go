package docker

import (
	"context"
	"fmt"
	"sort"

	"github.com/ErickMaria/envcontainer/internal/runtime/types"
)

func (docker *Docker) List(ctx context.Context, options map[string]types.ContainerOptions) error {

	paths := make([]string, 0, len(options))
	for key := range options {
		paths = append(paths, key)
	}
	sort.Strings(paths)

	fmt.Printf("%-20s %-20s\n", "NAME", "PATH")
	for _, path := range paths {
		fmt.Printf("%-20s %-20s\n", options[path].ContainerName, path)
	}

	return nil

}
