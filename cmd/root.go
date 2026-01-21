package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/ErickMaria/envcontainer/internal/runtime/docker"
	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/spf13/cobra"
)

type projectOptions struct {
	path      string
	getCloser bool
}

var (

	// # DOCKER API
	ctx       = context.Background()
	container = docker.NewDocker()

	rootCmd = &cobra.Command{
		Use:   "envcontainer",
		Short: "Create and manage reproducible development environments using Docker",
		Long:  "Envcontainer helps you create, run, and manage reproducible development environments backed by Docker containers. Use the subcommands to build images, start/stop containers, and list projects.",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	ops := projectOptions{}
	var err error
	ops.path, err = os.Getwd()
	if err != nil {
		slog.Error("%s", err.Error())
	}

	cobra.OnInitialize(initConfig)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(
		versionCommand(ops),
		initCommand(ops),
		listCommand(ops),
		buildCommand(ops),
		runCommand(ops),
		upCommand(ops),
		downCommand(ops),
	)
}

func initConfig() {
	err := template.Initialization()
	if err != nil {
		panic(err)
	}
}
