package main

import (
	"strings"

	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

func Down() cli.Command {
	return cli.Command{
		Desc: "remove all envcontainer configuration running in the current directory",
		Flags: cli.Flag{
			Values: map[string]cli.Values{
				"name": {
					Description: "container name",
				},
				"get-closer": {
					Defaulvalue: "true",
					Description: "will get the closest configuration file and remove all envcontainer",
				},
			},
		},
		Exec: func() {

			configFile, defaultMountDir, err := template.GetConfig(*cmd.Flags.Values["get-closer"].ValueBool)
			if err != nil {
				panic(err)
			}

			var containerName = configFile.Project.Name
			var noContainerNameSuffix = false

			if *cmd.Flags.Values["name"].ValueString != "" {
				containerName = *cmd.Flags.Values["name"].ValueString
				noContainerNameSuffix = true
			}

			commonLabels := map[string]string{
				"envcontainer/project-path": strings.TrimSuffix(defaultMountDir, ".envcontainer/"),
				"envcontainer/project-name": configFile.Project.Name,
				// "envcontainer/project-version":     configFile.Project.Version,
				// "envcontainer/project-description": configFile.Project.Description,
			}

			err = container.Down(ctx, types.ContainerOptions{
				ContainerName:     containerName,
				HostDirToBind:     path,
				NoContainerSuffix: noContainerNameSuffix,
				Networks:          configFile.Container.Networks,
				Labels:            commonLabels,
			})
			if err != nil {
				panic(err)
			}
		},
	}
}
