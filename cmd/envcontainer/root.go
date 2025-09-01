package main

import (
	"log"
	"os"

	"github.com/ErickMaria/envcontainer/internal/template"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

func Root() *cli.Command {

	var err error
	path, err = os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// # TEMPLATE FILE
	err = template.Initialization()
	if err != nil {
		panic(err)
	}

	cmd, cmds = cli.NewCommand(cli.CommandConfig{
		"build":    Build(),
		"up":       Up(),
		"down":     Down(),
		"run":      List(),
		"ls":       Run(),
		"template": Template(),
		"version":  Version(),
		"help":     Help(),
	})

	return cmd
}
