package command

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ErickMaria/envcontainer/internal/pkg/syscmd"
	"github.com/ErickMaria/envcontainer/internal/runtime"
	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

func StartCommand(ctx context.Context, configFile template.Envcontainer, errConfigFile error, container runtime.ContainerRuntime, flags cli.Flag) {

	// file := ""
	getCloser := *flags.Values["get-closer"].ValueBool
	if getCloser {
		file, err := syscmd.FindFileCloser(".envcontainer.yaml")
		if err != nil {
			panic(err)
		}

		if file != "" {
			configFile, err = template.UnmarshalWithFile(file)
			if err != nil {
				panic(err)
			}

		}
	} else if errConfigFile != nil {
		panic(errConfigFile)
	}

	// FIND USER
	if configFile.Container.User != "" {
		_, err := syscmd.FindUser(configFile.Container.User)
		if err != nil {
			panic(err)
		}
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
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

	autoStop := *flags.Values["auto-stop"].ValueBool

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
	})
	if err != nil {
		panic(err)
	}
}
