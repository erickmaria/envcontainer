package main

import (
	"fmt"
	"os"
	"strings"

	options "github.com/ErickMaria/envcontainer/pkg/cli"
	cmps "github.com/ErickMaria/envcontainer/internal/compose"
)

var cmd *options.Command
var cmds options.CommandConfig

func init() {

	dir, _ := os.Getwd()
	projectName := strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1]
	compose := cmps.Compose{}
	template := cmps.NewTemplate()

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

}

func main() {
	cmd.Listener()
}
