package main

import (
	"fmt"
	"os"
	"strings"

	options "github.com/ErickMaria/envcontainer/cli"
	"github.com/ErickMaria/envcontainer/common"
)

var cmd *options.Command
var cmds options.CommandConfig

func init() {

	dir, _ := os.Getwd()
	projectName := strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1]
	compose := common.Compose{}
	template := common.NewTemplate()

	cmd, cmds = options.NewCommand(options.CommandConfig{
		"init": options.Command{
			Flags: options.Flag{
				Values: map[string]options.Values{
					"build": options.Values{
						Defaulvalue: "false",
						Description: "init envcontainer and build configs",
					},
					"override": options.Values{
						Defaulvalue: "false",
						Description: "override envcontainer configs",
					},
				},
			},
			Quetion: options.Quetion{
				Queries: map[string]options.Query{
					"1_project": options.Query{
						Scene: "project_name [" + projectName + "]: ",
						Value: projectName,
					},
					"2_image": options.Query{
						Scene: "base_image [ubuntu:latest]: ",
						Value: "ubuntu:latest",
					},
					"3_ports": options.Query{
						Scene: "container_ports [\"80:80\"]: ",
					},
				},
			},
			Exec: func() {
				template.Init(cmd)
			},
			Desc: "create envcontainer template",
		},
		"build": options.Command{
			Desc: "prepare envcontainer to connect on container",
			Exec: func() {
				compose.Build()
			},
		},
		"up": options.Command{
			Flags: options.Flag{
				Values: map[string]options.Values{
					"shell": options.Values{
						Defaulvalue: "bash",
						Description: "terminal shell that must be used",
					},
				},
			},
			Desc: "creates the container and links to the project",
			Exec: func() {
				compose.Up(*cmd.Flags.Values["shell"].ValueString)
			},
		},
		"down": options.Command{
			Desc: "delete envcontainer container",
			Exec: func() {
				compose.Down()
			},
		},
		"delete": options.Command{
			Flags: options.Flag{
				Values: map[string]options.Values{
					"auto-approve": options.Values{
						Description: "skip confirmation",
						Defaulvalue: "false",
					},
				},
			},
			Exec: func() {
				template.Delete(cmd)
				compose.Delete()
			},
			Desc: "delete envcontainer configs",
		},
		"version": options.Command{
			Exec: func() {
				fmt.Println("0.2.0")
			},
			Desc: "show envcontainer version",
		},
		"help": options.Command{
			Exec: func() {
				options.Help(cmds)
			},
			Desc: "Run 'envcontainer COMMAND' for more information on a command. See: 'envcontainer help'",
		},
	})

	// cmd, cmds = options.NewCommand(options.CommandConfig{
	// 	options.INIT: options.Command{
	// 		Flags: options.Flag{
	// 			Command: options.INIT,
	// 			Values: map[string]options.Values{
	// 				"project": options.Values{
	// 					Defaulvalue: strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1],
	// 					Description: "project name",
	// 				},
	// 				"listener": options.Values{
	// 					Defaulvalue: "",
	// 					Description: "docker comtainer port listener",
	// 				},
	// 				"envfile": options.Values{
	// 					Defaulvalue: ".envcontainer/compose/.env",
	// 					Description: "docker environemt file",
	// 				},
	// 				"image": options.Values{
	// 					Defaulvalue: "ubuntu",
	// 					Description: "dockerfile image",
	// 				},
	// 				"no-build": options.Values{
	// 					Defaulvalue: "false",
	// 					Description: "init envcontainer and build configs",
	// 				},
	// 				"override": options.Values{
	// 					Defaulvalue: "false",
	// 					Description: "override envcontainer configs",
	// 				},
	// 			},
	// 		},
	// 		Exec: func() {
	// 			options.Init(cmd.Flags)
	// 		},
	// 		Desc: "create envcontainer blueprint",
	// 	},
	// 	options.BUILD: options.Command{
	// 		Flags: options.Flag{
	// 			Command: options.BUILD,
	// 		},
	// 		Desc: "prepare envcontainer to connect on container",
	// 		Exec: func() {
	// 			options.Build()
	// 		},
	// 	},
	// 	options.START: options.Command{
	// 		Flags: options.Flag{
	// 			Command: options.START,
	// 			Values: map[string]options.Values{
	// 				"shell": options.Values{
	// 					Defaulvalue: "bash",
	// 					Description: "terminal shell that must be used",
	// 				},
	// 			},
	// 		},
	// 		Desc: "creates the container and links to the project",
	// 		Exec: func() {
	// 			options.Start(cmd.Flags)
	// 		},
	// 	},
	// 	options.STOP: options.Command{
	// 		Flags: options.Flag{
	// 			Command: options.STOP,
	// 		},
	// 		Desc: "delete envcontainer container",
	// 		Exec: func() {
	// 			options.Stop()
	// 		},
	// 	},
	// 	options.DELETE: options.Command{
	// 		Flags: options.Flag{
	// 			Command: options.DELETE,
	// 			Values: map[string]options.Values{
	// 				"auto-approve": options.Values{
	// 					Description: "skip confirmation",
	// 					Defaulvalue: "false",
	// 				},
	// 			},
	// 		},
	// 		Exec: func() {
	// 			options.Delete(cmd.Flags)
	// 		},
	// 		Desc: "delete envcontainer configs",
	// 	},
	// 	options.VERSION: options.Command{
	// 		Exec: func() {
	// 			fmt.Println("0.0.1-alpha")
	// 		},
	// 		Desc: "show envcontainer version",
	// 	},
	// 	options.HELP: options.Command{
	// 		Exec: func() {
	// 			options.Help(cmds)
	// 		},
	// 		Desc: "Run 'envcontainer COMMAND' for more information on a command. See: 'envcontainer help'",
	// 	},
	// })

}

func main() {
	cmd.Listener()
}
