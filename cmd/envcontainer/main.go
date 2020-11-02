package main

import (
	"fmt"
	"os"

	"github.com/ErickMaria/envcontainer/internal/options"
)

var cmds map[string]options.Command

func init() {

	cmds = map[string]options.Command{
		options.INIT: options.Command{
			Flags: options.NewFlag(options.Flag{
				Command: "init",
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
			}),
			Exec: func() {
				options.Init(cmds[options.INIT].Flags)
			},
			Desc: "create envcontainer blueprint",
		},

		options.HELP: options.Command{
			Exec: func() {
				options.Help(cmds)
			},
			Desc: "Run 'envcontainer COMMAND --help' for more information on a command. See: 'envcontainer help'",
		},
	}

}

func main() {

	switch os.Args[1] {
	case options.INIT:
		cmds[options.INIT].Exec()
	case options.HELP:
		cmds[options.HELP].Exec()
	default:
		fmt.Printf("envcontainer: '%s' is not a envcontainer command.\n", os.Args[1])
	}

}
