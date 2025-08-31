package main

import (
	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

func List() cli.Command {
	return cli.Command{
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
	}
}
