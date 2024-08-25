package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/runtime/docker"
	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

var cmd *cli.Command
var cmds cli.CommandConfig

func init() {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// # TEMPLATE FILE
	err = template.Initialization()
	if err != nil {
		panic(err)
	}

	// # DOCKER API
	ctx := context.Background()
	container := docker.NewDocker()

	// CLI
	cmd, cmds = cli.NewCommand(cli.CommandConfig{
		"build": cli.Command{
			Desc: "build a image using envcontainer configuration in the current directory",
			Exec: func() {

				configFile, _, err := template.GetConfig(false)
				if err != nil {
					panic(err)
				}

				err = container.Build(ctx, types.BuildOptions{
					ImageName:    configFile.Project.Name,
					Dockerfile:   configFile.Container.Build,
					BuildContext: configFile.GetTmpDockerfileDir(),
				})
				if err != nil {
					panic(err)
				}
			},
		},
		"up": cli.Command{
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
		},
		"down": cli.Command{
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

				configFile, _, err := template.GetConfig(*cmd.Flags.Values["get-closer"].ValueBool)
				if err != nil {
					panic(err)
				}

				var containerName = configFile.Project.Name
				var noContainerNameSuffix = false

				if *cmd.Flags.Values["name"].ValueString != "" {
					containerName = *cmd.Flags.Values["name"].ValueString
					noContainerNameSuffix = true
				}

				err = container.Down(ctx, types.ContainerOptions{
					ContainerName:     containerName,
					HostDirToBind:     path,
					NoContainerSuffix: noContainerNameSuffix,
				})
				if err != nil {
					panic(err)
				}
			},
		},
		"run": cli.Command{
			Desc: "execute an .envcontainer on the current directory without saving it locally",
			Flags: cli.Flag{
				Values: map[string]cli.Values{
					"name": {
						Description: "container name",
					},
					"image": {
						Description: "envcontainer image",
					},
					"command": {
						Description: "execute command inside container",
					},
				},
			},
			Exec: func() {

				name := "envcontainer"
				image := *cmd.Flags.Values["image"].ValueString
				command := *cmd.Flags.Values["command"].ValueString

				err := container.Run(ctx, types.ContainerOptions{
					ContainerName: name,
					ImageName:     image,
					Commands:      strings.Split(strings.Trim(command, " "), " "),
					AutoStop:      true,
					HostDirToBind: path,
				})
				if err != nil {
					panic(err)
				}

			},
		},
		"ls": cli.Command{
			Desc: "list envcontainers",
			Exec: func() {
				configFiles, err := template.List()
				if err != nil {
					panic(err)
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

				err = container.List(ctx, containerOpts)
				if err != nil {
					panic(err)
				}

			},
		},
		"version": cli.Command{
			Exec: func() {
				fmt.Println("Version: 2.0.0")
			},
			Desc: "show envcontainer version",
		},
		"help": cli.Command{
			Exec: func() {
				cli.Help(cmds)
			},
			Desc: "Run " + cli.ExecutableName() + " COMMAND' for more information on a command. See: '" + cli.ExecutableName() + " help'",
		},
	})

}

func main() {
	cmd.Listener()
}
