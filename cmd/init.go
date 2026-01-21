package cmd

import (
	"fmt"

	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/spf13/cobra"
)

type initOptions struct{}

func initCommand(projOpts projectOptions) *cobra.Command {
	ops := initOptions{}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a starter .envcontainer.yaml configuration file",
		Long:  "Generate a starter .envcontainer.yaml in the current project directory following the Envcontainer schema.",
		Run: func(cmd *cobra.Command, args []string) {
			ops.execute(projOpts)
		},
	}

	return cmd
}

func (i initOptions) execute(projOpts projectOptions) {

	err := template.NewConfigFile(projOpts.path)
	if err != nil {
		fmt.Printf("Error creating .envcontainer.yaml: %s\n", err.Error())
		return
	}
	dest := fmt.Sprintf("%s/.envcontainer.yaml", projOpts.path)

	fmt.Printf("Created %s\n", dest)
}
