package cmd

import (
	"strings"

	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/spf13/cobra"
)

type downOptions struct {
	containerName string
}

func downCommand(projOpts projectOptions) *cobra.Command {

	ops := downOptions{}
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Stop and remove containers created by envcontainer in the current directory",
		Long:  "Stop and remove containers and related resources created by envcontainer for the current project. Use --name to target a specific container name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ops.execute(projOpts)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&ops.containerName, "name", "n", "envcontainer", "Container name (defaults to the project name unless --name is provided)")
	flags.BoolVarP(&projOpts.getCloser, "get-closer", "", false, "Search parent dirs and use the nearest envcontainer configuration file")

	return cmd
}

func (d downOptions) execute(projOpts projectOptions) error {

	configFile, defaultMountDir, err := template.GetConfig(projOpts.getCloser)
	if err != nil {
		panic(err)
	}

	var containerName = configFile.Project.Name
	var noContainerNameSuffix = false

	if d.containerName != "" {
		containerName = d.containerName
		noContainerNameSuffix = true
	}

	commonLabels := map[string]string{
		"envcontainer/project-path": strings.TrimSuffix(defaultMountDir, ".envcontainer/"),
		"envcontainer/project-name": configFile.Project.Name,
		// "envcontainer/project-version":     configFile.Project.Version,
		// "envcontainer/project-description": configFile.Project.Description,
	}

	return container.Down(ctx, types.ContainerOptions{
		ContainerName:     containerName,
		HostDirToBind:     projOpts.path,
		NoContainerSuffix: noContainerNameSuffix,
		Networks:          configFile.Container.Networks,
		Labels:            commonLabels,
	})
}
