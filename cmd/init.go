package cmd

import (
	"fmt"

	"github.com/ErickMaria/envcontainer/internal/template"
	tplTypes "github.com/ErickMaria/envcontainer/internal/template/types"
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
	// update cache by adding the new config directly (avoid full rescan)
	updateCacheWithConfig(dest)
}

func updateCacheWithConfig(configPath string) {
	env, err := template.UnmarshalWithFile(configPath)
	if err != nil {
		fmt.Printf("Warning: created config but failed to unmarshal for cache: %v\n", err)
		return
	}

	cache, err := template.ListCached()
	if err != nil {
		fmt.Printf("Warning: failed to read cache: %v\n", err)
		return
	}

	if cache == nil {
		cache = map[string]tplTypes.Envcontainer{}
	}
	cache[configPath] = env

	if err := template.SaveCache(cache); err != nil {
		fmt.Printf("Warning: failed to update cache: %v\n", err)
	}
}
