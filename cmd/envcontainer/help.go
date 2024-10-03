package main

import (
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

func Help() cli.Command {
	return cli.Command{
		Exec: func() {
			cli.Help(cmds)
		},
		Desc: "Run " + cli.ExecutableName() + " COMMAND' for more information on a command. See: '" + cli.ExecutableName() + " help'",
	}
}
