package cmd

import (
	"strings"

	"github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/spf13/cobra"
)

type runOptions struct {
	name    string
	image   string
	command string
}

func runCommand(projOpts projectOptions) *cobra.Command {
	ops := runOptions{}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a one-off container from an image without saving configuration",
		Long:  "Run a transient container from a specified image and execute a command inside it. The container is not persisted as an envcontainer project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ops.execute(projOpts)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&ops.name, "name", "n", "envcontainer", "Container name (default: envcontainer)")
	flags.StringVarP(&ops.image, "image", "i", "", "Container image to run (required)")
	cmd.MarkFlagRequired("image")
	flags.StringVarP(&ops.command, "command", "c", "", "Command to execute inside the container (required). For multiple args, wrap in quotes")
	cmd.MarkFlagRequired("command")

	return cmd
}

func (r runOptions) execute(projOpts projectOptions) error {

	return container.Run(ctx, types.ContainerOptions{
		ContainerName: r.name,
		ImageName:     r.image,
		Commands:      strings.Split(strings.Trim(r.command, " "), " "),
		AutoStop:      true,
		HostDirToBind: projOpts.path,
	})
}
