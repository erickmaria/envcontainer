package main

import (
	"fmt"
	"os"
	"strings"

	cmps "github.com/ErickMaria/envcontainer/internal/compose"
	"github.com/ErickMaria/envcontainer/internal/envasync"
	"github.com/ErickMaria/envcontainer/internal/envconfig"
	"github.com/ErickMaria/envcontainer/pkg/cli"
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
	upasync := envasync.UpAsync{}

	cmd, cmds = options.NewCommand(options.CommandConfig{
		"init": options.Command{
			Flags: options.Flag{
				Values: map[string]options.Values{
					"build": options.Values{
						Defaulvalue: "false",
						Description: "build a image using envcontainer configuration",
					},
					"override": options.Values{
						Defaulvalue: "false",
						Description: "override envcontainer configuration",
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
			Desc: "initialize the default template in the current directory",
		},
		"build": options.Command{
			Desc: "build a image using envcontainer configuration in the current directory",
			Exec: func() {
				compose.Build()
			},
		},
		"run": options.Command{
			Flags: options.Flag{
				Values: map[string]options.Values{
					"shell": options.Values{
						Defaulvalue: "bash",
						Description: "terminal shell that must be used",
					},
				},
			},
			Desc: "run the envcontainer configuration to start the container and link it to the current directory",
			Exec: func() {
				compose.Up(*cmd.Flags.Values["shell"].ValueString)
			},
		},
		"stop": options.Command{
			Desc: "stop all envcontainer configuration running in the current directory",
			Exec: func() {
				compose.Down()
			},
		},
		"save": options.Command{
			Exec: func() {
				config.Save()
			},
			Desc: "save your local .envcontainer directory",
		},
		"list": options.Command{
			Exec: func() {
				config.List()
			},
			Desc: "list all your .envcontainer directory saved",
		},
		"get": options.Command{
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
			Desc: "get .envcontainer and put in current directory",
		},
		"exec": options.Command{
			Flags: options.Flag{
				Values: map[string]options.Values{
					"name": options.Values{
						Description: "envcontainer configuration name",
					},
				},
			},
			Exec: func() {
				upasync.Start(cmd)
			},
			Desc: "execute an .envcontainer on the current directory without saving it locally",
		},
		// "remove": options.Command{
		// 	Flags: options.Flag{
		// 		Values: map[string]options.Values{
		// 			"auto-approve": options.Values{
		// 				Description: "skip confirmation",
		// 				Defaulvalue: "false",
		// 			},
		// 		},
		// 	},
		// 	RunBeforeAll: func() {
		// 		template.Delete(cmd)
		// 	},
		// 	Exec: func() {
		// 		compose.Delete()
		// 	},
		// 	Desc: "delete local .envcontainer",
		// },
		"version": options.Command{
			Exec: func() {
				fmt.Println("Version: 0.5.0")
			},
			Desc: "show envcontainer version",
		},
		"help": options.Command{
			Exec: func() {
				options.Help(cmds)
			},
			Desc: "Run " + cli.ExecutableName() + " COMMAND' for more information on a command. See: '" + cli.ExecutableName() + " help'",
		},
	})

}

func main() {
	cmd.Listener()
}
