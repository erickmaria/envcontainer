package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/pkg/syscmd"
	"github.com/ErickMaria/envcontainer/internal/runtime/docker"
	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

var cmd *cli.Command
var cmds cli.CommandConfig

func getConfig(getCloser bool) template.Envcontainer {

	configFile, errConfigFile := template.Unmarshal()

	if getCloser {
		file, err := syscmd.FindFileCloser(".envcontainer.yaml")
		if err != nil {
			panic(err)
		}

		pwd, _ := os.Getwd()
		for i := 0; i < strings.Count(file, "../"); i++ {
			pwd = strings.Join(strings.Split(pwd, "/")[:len(strings.Split(pwd, "/"))-1], "/")

		}

		if file != "" {
			configFile, err = template.UnmarshalWithFile(file)
			if err != nil {
				panic(err)
			}

		}

		configFile.SetMountDir(pwd + "/.envcontainer/")
		configFile.BuildMount()

	} else if errConfigFile != nil {
		panic(errConfigFile)
	}
	return configFile

}

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

				configFile := getConfig(false)

				// FIND USER
				if configFile.Container.User != "" {
					_, err := syscmd.FindUser(configFile.Container.User)
					if err != nil {
						panic(err)
					}
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
		"start": cli.Command{
			Flags: cli.Flag{
				Values: map[string]cli.Values{
					"get-closer": {
						Defaulvalue: "true",
						Description: "will stop current container running and get the closest config file to run a new container",
					},
					"auto-stop": {
						Defaulvalue: "false",
						Description: "terminal shell that must be used",
					},
				},
			},
			Desc: "run the envcontainer configuration to start the container and link it to the current directory",
			Exec: func() {

				configFile := getConfig(*cmd.Flags.Values["get-closer"].ValueBool)

				// FIND USER
				if configFile.Container.User != "" {
					_, err := syscmd.FindUser(configFile.Container.User)
					if err != nil {
						panic(err)
					}
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

				if configFile.AutoStop {
					autoStop = configFile.AutoStop
				}

				fmt.Println("Stating container...")

				err = container.Start(ctx, types.ContainerOptions{
					AutoStop:        autoStop,
					ContainerName:   configFile.Project.Name,
					Ports:           configFile.Container.Ports,
					PullImageAlways: false,
					User:            configFile.Container.User,
					HostDirToBind:   path,
					Mounts:          configFile.Mounts,
				})
				if err != nil {
					panic(err)
				}
			},
		},
		"stop": cli.Command{
			Desc: "stop all envcontainer configuration running in the current directory",
			Flags: cli.Flag{
				Values: map[string]cli.Values{
					"name": {
						Description: "container name",
					},
					"get-closer": {
						Defaulvalue: "true",
						Description: "will stop current container running and get the closest config file to run a new container",
					},
				},
			},
			Exec: func() {

				configFile := getConfig(*cmd.Flags.Values["get-closer"].ValueBool)

				var containerName = configFile.Project.Name
				var noContainerNameSuffix = false

				if *cmd.Flags.Values["name"].ValueString != "" {
					containerName = *cmd.Flags.Values["name"].ValueString
					noContainerNameSuffix = true
				}

				err := container.Stop(ctx, types.ContainerOptions{
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
			Desc: "execute an .envcontainer on the current directory without saving it locally",
		},
		"version": cli.Command{
			Exec: func() {
				fmt.Println("Version: 1.0.0-beta5")
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
