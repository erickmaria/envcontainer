package main

import (
	"github.com/ErickMaria/envcontainer/internal/options"
)

var cmd *options.Command
var cmds options.CommandConfig

func init() {

	cmd, cmds = options.NewCommand(options.CommandConfig{
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
				options.Init(cmd.Flags)
			},
			Desc: "create envcontainer blueprint",
		},
		options.RUN: options.Command{
			Flags: options.Flag{
				Command: options.RUN,
			},
			Desc: "creates the container and links to the project",
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
				options.Delete(cmd.Flags)
			},
			Desc: "delete envcontainer configs",
		},
		options.HELP: options.Command{
			Exec: func() {
				options.Help(cmds)
			},
			Desc: "Run 'envcontainer COMMAND' for more information on a command. See: 'envcontainer help'",
		},
	})

}

func main() {
	cmd.Listener()
}
