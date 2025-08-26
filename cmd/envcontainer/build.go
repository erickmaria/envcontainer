package main

import (
	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

func Build() cli.Command {
	return cli.Command{
		Desc: "build a image using envcontainer configuration in the current directory",
		Exec: func() {

			configFile, _, err := template.GetConfig(false)
			if err != nil {
				panic(err)
			}

			if configFile.Container.NetworkMode == "" {
				configFile.Container.NetworkMode = "default"
			}

			err = container.Build(ctx, types.BuildOptions{
				ImageName:    configFile.Project.Name,
				Dockerfile:   configFile.Container.Build,
				BuildContext: configFile.GetTmpDockerfileDir(),
				NetworkMode:  configFile.Container.NetworkMode,
			})
			if err != nil {
				panic(err)
			}
		},
	}
}
