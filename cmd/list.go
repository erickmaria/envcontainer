package cmd

import (
	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	tplTypes "github.com/ErickMaria/envcontainer/internal/template/types"
	"github.com/spf13/cobra"
)

type listOptions struct {
	Refresh bool
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

	flags := cmd.Flags()
	flags.BoolVarP(&ops.Refresh, "refresh", "f", false, "refresh the list of envcontainer projects by re-scanning the filesystem for configuration files")

	return cmd
}

func (l listOptions) execute(projOpts projectOptions) error {
	var configFiles map[string]tplTypes.Envcontainer
	var err error

	if l.Refresh {
		configFiles, err = template.RefreshCache()
		if err != nil {
			return err
		}
	} else {
		configFiles, err = template.ListCached()
		if err != nil {
			return err
		}
		// if cache is empty, populate it by scanning
		if len(configFiles) == 0 {
			configFiles, err = template.RefreshCache()
			if err != nil {
				return err
			}
		}
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
