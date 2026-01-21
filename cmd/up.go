package cmd

import (
	"fmt"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/spf13/cobra"

	"github.com/ErickMaria/envcontainer/internal/template"
)

type upOptions struct {
	autoStop bool
	code     bool
	cursor   bool
	host     string
	port     uint32
}

func upCommand(projOpts projectOptions) *cobra.Command {
	ops := upOptions{}
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Start a container from envcontainer config and bind it to the current directory",
		Long:  "Start (or recreate) a container defined by the envcontainer configuration and bind the current directory into the container. Use flags to control editor integration and SSH settings.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ops.execute(projOpts)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&projOpts.getCloser, "get-closer", "", false, "Search parent dirs and use the nearest envcontainer configuration file")
	flags.BoolVarP(&ops.autoStop, "auto-stop", "", false, "Automatically stop the container when the main process exits")
	flags.BoolVarP(&ops.code, "code", "c", false, "Open the project in VS Code via SSH")
	flags.BoolVarP(&ops.cursor, "cursor", "", false, "Open the project in Cursor editor via SSH")
	flags.StringVarP(&ops.host, "host", "a", "", "SSH host address used by the editor remote (only used with --code or --cursor)")
	flags.Uint32VarP(&ops.port, "port", "p", 22, "SSH port used by the editor remote (only used with --code or --cursor)")

	return cmd
}

func (u upOptions) execute(projOpts projectOptions) error {

	if u.code && u.cursor {
		return fmt.Errorf("cannot use both --code and --cursor flags together. Please use only one editor at a time")
	}

	configFile, defaultMountDir, err := template.GetConfig(projOpts.getCloser)

	if err != nil {
		panic(err)
	}

	if configFile.AlwaysUpdate {

		err := container.AlwaysUpdate(ctx, types.BuildOptions{
			ImageName:  configFile.Project.Name,
			Dockerfile: configFile.Container.Build,
		})
		if err != nil {
			panic(err)
		}
	}

	commonLabels := map[string]string{
		"envcontainer/project-path":        strings.TrimSuffix(defaultMountDir, ".envcontainer/"),
		"envcontainer/project-name":        configFile.Project.Name,
		"envcontainer/project-version":     configFile.Project.Version,
		"envcontainer/project-description": configFile.Project.Description,
	}

	if configFile.Container.NetworkMode == "" {
		configFile.Container.NetworkMode = "default"
	}

	var editor string
	if u.code {
		editor = "code"
	} else if u.cursor {
		editor = "cursor"
	}

	return container.Up(ctx, types.ContainerOptions{
		AutoStop:        u.autoStop,
		ContainerName:   configFile.Project.Name,
		Ports:           configFile.Container.Ports,
		PullImageAlways: false,
		Shell:           configFile.Container.Shell,
		HostDirToBind:   projOpts.path,
		Mounts:          configFile.Mounts,
		DefaultMountDir: defaultMountDir,
		NetworkMode:     configFile.Container.NetworkMode,
		Networks:        configFile.Container.Networks,
		Labels:          commonLabels,
	}, editor, u.host, u.port)
}
