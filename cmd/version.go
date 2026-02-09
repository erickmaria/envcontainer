package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type versionOptions struct {
}

func versionCommand(projOpts projectOptions) *cobra.Command {
	ops := versionOptions{}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show envcontainer CLI version",
		Long:  "Print the currently installed envcontainer CLI version.",
		Run: func(cmd *cobra.Command, args []string) {
			ops.execute(projOpts)
		},
	}

	return cmd
}

func (v versionOptions) execute(projOpts projectOptions) {
	fmt.Println("Version: 2.10.0")
}
