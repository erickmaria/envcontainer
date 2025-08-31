package main

import (
	"strings"

	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

func Run() cli.Command {

	return cli.Command{
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
	}
}
