package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ErickMaria/envcontainer/internal/pkg/randon"
	"github.com/ErickMaria/envcontainer/internal/pkg/syscmd"
	"github.com/ErickMaria/envcontainer/internal/runtime/docker"
	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

var cmd *cli.Command
var cmds cli.CommandConfig

func init() {

	// RANDON SEED
	rand.Seed(time.Now().UnixNano())

	// # TEMPLATE FILE
	err := template.Initialization()
	if err != nil {
		panic(err)
	}

	configFile, err := template.Unmarshal()
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

				// FIND USER
				if configFile.Container.User != "" {
					_, err := syscmd.FindUser(configFile.Container.User)
					if err != nil {
						panic(err)
					}
				}

				err = container.Build(ctx, types.BuildOptions{
					ImageName:  configFile.Project.Name,
					Dockerfile: configFile.Container.Build,
				})
				if err != nil {
					panic(err)
				}
			},
		},
		"start": cli.Command{
			Flags: cli.Flag{
				Values: map[string]cli.Values{
					"auto-stop": {
						Defaulvalue: "false",
						Description: "terminal shell that must be used",
					},
				},
			},
			Desc: "run the envcontainer configuration to start the container and link it to the current directory",
			Exec: func() {

				// FIND USER
				// FIND USER
				if configFile.Container.User != "" {
					_, err := syscmd.FindUser(configFile.Container.User)
					if err != nil {
						panic(err)
					}
				}
				autoStop := *cmd.Flags.Values["auto-stop"].ValueBool

				err = container.Start(ctx, types.ContainerOptions{
					AutoStop:        autoStop,
					ContainerName:   configFile.Project.Name,
					Ports:           configFile.Container.Ports,
					PullImageAlways: false,
					User:            configFile.Container.User,
				})
				if err != nil {
					panic(err)
				}
			},
		},
		"stop": cli.Command{
			Desc: "stop all envcontainer configuration running in the current directory",
			Exec: func() {
				err := container.Stop(ctx, configFile.Project.Name)
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

				name := "envcontainer_" + randon.RandStringRunes(5)
				image := *cmd.Flags.Values["image"].ValueString
				command := *cmd.Flags.Values["command"].ValueString

				err := container.Run(ctx, types.ContainerOptions{
					ContainerName: name,
					ImageName:     image,
					Commands:      strings.Split(strings.Trim(command, " "), " "),
					AutoStop:      true,
				})
				if err != nil {
					panic(err)
				}

			},
			Desc: "execute an .envcontainer on the current directory without saving it locally",
		},
		"version": cli.Command{
			Exec: func() {
				fmt.Println("Version: 0.5.0")
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
