package main

import (
	"fmt"
	"os"

	"github.com/ErickMaria/envcontainer/internal/options"
)

var cmds = options.CommandConfig

func init() {

	cmds = map[string]options.Command{
		options.INIT: options.Command{
			Flags: options.Flag{
				Command: options.INIT,
				Values: map[string]options.Values{
					"project": options.Values{
						Defaulvalue: "app",
						Description: "project name",
					},
					"listener": options.Values{
						Defaulvalue: "",
						Description: "docker comtainer port listener",
					},
					"envfile": options.Values{
						Defaulvalue: "env/.variables",
						Description: "docker environemt file",
					},
					"image": options.Values{
						Defaulvalue: "ubuntu",
						Description: "dockerfile image",
					},
				},
			},
			Exec: func() {
				options.Init(cmds[options.INIT].Flags)
			},
			Desc: "create envcontainer blueprint",
		},
		options.RUN: options.Command{
			Flags: options.Flag{
				Command: options.RUN,
			},
			Desc: "",
			Exec: func() {
				options.Run()
			},
		},
		options.DELETE: options.Command{
			Flags: options.Flag{
				Command: options.DELETE,
				Values: map[string]options.Values{
					"auto-approve": options.Values{
						Description: "skip confirmation (yes/no)",
					},
				},
			},
			Exec: func() {
				options.Delete(cmds[options.DELETE].Flags)
			},
			Desc: "delete envcontainer configs",
		},
		options.HELP: options.Command{
			Exec: func() {
				options.Help(cmds)
			},
			Desc: "Run 'envcontainer COMMAND' for more information on a command. See: 'envcontainer help'",
		},
	}

}

func main() {

	flgs := cmds[os.Args[1]].Flags
	flgs.Register()

	switch os.Args[1] {
	case options.INIT:
		cmds[options.INIT].Exec()
	case options.RUN:
		cmds[options.RUN].Exec()
	case options.DELETE:
		cmds[options.DELETE].Exec()
	case options.HELP:
		cmds[options.HELP].Exec()
	default:
		fmt.Printf("envcontainer: '%s' is not a envcontainer command.\n", os.Args[1])
	}

}
