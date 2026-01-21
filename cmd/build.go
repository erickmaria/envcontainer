package cmd

import (
	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/spf13/cobra"
)

type buildOptions struct {
}

func buildCommand(projOpts projectOptions) *cobra.Command {
	ops := buildOptions{}
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build a Docker image from the envcontainer configuration in the current directory",
		Long:  "Builds a Docker image using the envcontainer configuration found in the current directory. Use --get-closer to search parent directories for a configuration file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ops.execute(projOpts)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&projOpts.getCloser, "get-closer", "", false, "Search parent directories and use the nearest envcontainer configuration file")

	return cmd
}

func (b buildOptions) execute(projOpts projectOptions) error {

	configFile, _, err := template.GetConfig(projOpts.getCloser)
	if err != nil {
		return err
	}

	if configFile.Container.NetworkMode == "" {
		configFile.Container.NetworkMode = "default"
	}

	return container.Build(ctx, types.BuildOptions{
		ImageName:    configFile.Project.Name,
		Dockerfile:   configFile.Container.Build,
		BuildContext: template.GetTmpDockerfileDir(configFile),
		NetworkMode:  configFile.Container.NetworkMode,
	})
}
