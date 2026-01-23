package cmd

import (
	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/spf13/cobra"
)

type listOptions struct {
}

func listCommand(projOpts projectOptions) *cobra.Command {
	ops := listOptions{}
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List discovered envcontainer projects and their status",
		Long:    "Discover envcontainer projects on the system and show basic information and current container status. Use this to find available projects and their runtime state.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ops.execute(projOpts)
		},
	}

	return cmd
}

func (l listOptions) execute(projOpts projectOptions) error {
	configFiles, err := template.List()
	if err != nil {
		return err
	}

	var containerOpts = map[string]types.ContainerOptions{}
	for path, configs := range configFiles {
		containerOpts[path] = types.ContainerOptions{
			ContainerName: configs.Project.Name,
			Ports:         configs.Container.Ports,
			Shell:         configs.Container.Shell,
			HostDirToBind: path,
			Mounts:        configs.Mounts,
		}
	}

	return container.List(ctx, containerOpts)
}
