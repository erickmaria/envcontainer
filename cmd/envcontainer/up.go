package main

import (
	"fmt"

	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

func Up() cli.Command {
	return cli.Command{
		Flags: cli.Flag{
			Values: map[string]cli.Values{
				"get-closer": {
					Defaulvalue: "true",
					Description: "will get the closest configuration file to run a new container",
				},
				"auto-stop": {
					Defaulvalue: "false",
					Description: "terminal shell that must be used",
				},
				"code": {
					Defaulvalue: "false",
					Description: "open with vscode",
				},
			},
		},
		Desc: "run the envcontainer configuration to start the container and link it to the current directory",
		Exec: func() {

			configFile, defaultMountDir, err := template.GetConfig(*cmd.Flags.Values["get-closer"].ValueBool)
			if err != nil {
				panic(err)
			}

			if configFile.AlwaysUpdate {
				fmt.Println("Restat container...")

				err := container.AlwaysUpdate(ctx, types.BuildOptions{
					ImageName:  configFile.Project.Name,
					Dockerfile: configFile.Container.Build,
				})
				if err != nil {
					panic(err)
				}
			}

			autoStop := *cmd.Flags.Values["auto-stop"].ValueBool
			code := *cmd.Flags.Values["code"].ValueBool

			if configFile.AutoStop {
				autoStop = configFile.AutoStop
			}

			fmt.Println("Stating container...")

			err = container.Up(ctx, types.ContainerOptions{
				AutoStop:        autoStop,
				ContainerName:   configFile.Project.Name,
				Ports:           configFile.Container.Ports,
				PullImageAlways: false,
				Shell:           configFile.Container.Shell,
				HostDirToBind:   path,
				Mounts:          configFile.Mounts,
				DefaultMountDir: defaultMountDir,
			}, code)
			if err != nil {
				panic(err)
			}
		},
	}
}
