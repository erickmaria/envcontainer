package compose

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/pkg/handler/errors"
	"github.com/ErickMaria/envcontainer/internal/pkg/syscmd"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

type Template struct {
	Home        string
	PathDefault string
	Dockerfile  string
	Compose     string
	Env         string
}

func NewTemplate() Template {
	return Template{
		Home:        HOME,
		PathDefault: PATH_COMPLETE,
		Dockerfile:  ".envcontainer/Dockerfile",
		Compose:     COMPOSE,
		Env:         ".envcontainer/compose/.env",
	}
}

func (template *Template) Init(commad *cli.Command) {

	queries := commad.Quetion.Queries
	values := commad.Flags.Values

	if !*values["override"].ValueBool {
		if _, err := os.Stat(template.Home); !os.IsNotExist(err) {
			fmt.Println("\nenvcontainer has already started in this folder.\nif you want to override use: --override flag with init comand")
			os.Exit(0)
		}
	}

	var ports = []string{}
	if queries["3_ports"].Value != "" {
		ports = append(ports, queries["3_ports"].Value)
	}

	compose := Compose{
		Version: "3.6",
		Services: Services{
			Environment: Environment{
				ContainerName: strings.ToLower(queries["1_project"].Value),
				Volumes: []Volumes{
					Volumes{
						Type:   "bind",
						Source: "../../",
						Target: "/home/envcontainer/" + queries["1_project"].Value,
					},
				},
				Build: Build{
					Dockerfile: "Dockerfile",
					Context:    "../",
				},
				Image:      strings.ToLower("envcontainer/" + queries["1_project"].Value),
				Ports:      ports,
				WorkingDir: "/home/envcontainer/" + queries["1_project"].Value,
				EnvFile: []string{
					"$PWD/" + template.Env,
				},
				Privileged: true,
				StdinOpen:  true,
				Tty:        true,
			},
		},
	}

	data := compose.Marshal()

	err := syscmd.CreatePath(template.PathDefault)
	errors.Throw("envcontainer: error to create folders, check permissions", err)

	syscmd.CreateFile(template.Dockerfile, []byte("FROM "+queries["2_image"].Value))
	syscmd.CreateFile(template.Compose, data)

	envContent := fmt.Sprintf("COMPOSE_PROJECT_NAME=%s\nCOMPOSE_IGNORE_ORPHANS=True", queries["1_project"].Value)
	syscmd.CreateFile(template.Env, []byte(envContent))

	if *values["build"].ValueBool {
		compose.Build()
	}

	fmt.Println("\nenvcontainer initialized!")
}

func (template *Template) Delete(commad *cli.Command) {

	autoApprove := commad.Flags.Values["auto-approve"]
	if !*autoApprove.ValueBool {
		fmt.Print("do you're have sure? (yes/no): ")
		reader := bufio.NewReader(os.Stdin)
		confirmation, _, err := reader.ReadLine()
		if err != nil {
			errors.Throw("", err)
		}

		v := string(confirmation)

		switch strings.ToLower(v) {
		case "yes":
			break
		case "no":
			os.Exit(0)
			return
		default:
			fmt.Println("envcontainer: invalid value")
			os.Exit(0)
		}
	}
}

func (template *Template) CheckEnvcontainerExists(flag *cli.Flag) {

	override := flag.Values["override"].ValueBool

	if syscmd.ExistsPath(template.Home) {

		if !*override {
			fmt.Printf("%s already exists, use --override\n", template.Home)
			os.Exit(0)
		}
	}
}
