package main

import (
	"fmt"
	"os"
	"strings"

	cmps "github.com/ErickMaria/envcontainer/internal/compose"
	"github.com/ErickMaria/envcontainer/internal/envconfig"
	options "github.com/ErickMaria/envcontainer/pkg/cli"
)

var cmd *options.Command
var cmds options.CommandConfig

func init() {

	envconfig.CreateIfNotExists()

	dir, _ := os.Getwd()
	projectName := strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1]
	compose := cmps.Compose{}
	template := cmps.NewTemplate()
	config := envconfig.Config{}

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
			RunBeforeAll: func() {
				template.CheckEnvcontainerExists(&cmd.Flags)
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
		"config-save": options.Command{
			Exec: func() {
				config.Save()
			},
		},
		"config-list": options.Command{
			Exec: func() {
				config.List()
			},
		},
		"config-get": options.Command{
			Flags: options.Flag{
				Values: map[string]options.Values{
					"name": options.Values{
						Description: "envcontainer configuration name",
					},
				},
			},
			Exec: func() {
				config.Get(cmd)
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
			RunBeforeAll: func() {
				template.Delete(cmd)
			},
			Exec: func() {
				compose.Delete()
			},
			Desc: "delete envcontainer configs",
		},
		"version": options.Command{
			Exec: func() {
				fmt.Println("0.4.0")
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
