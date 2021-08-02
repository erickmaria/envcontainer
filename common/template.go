package common

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ErickMaria/envcontainer/cli"
)

var (
	COMPOSE string = ".envcontainer/compose/docker-compose.yaml"
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
		Home:        ".envcontainer",
		PathDefault: ".envcontainer/compose",
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
				ContainerName: queries["1_project"].Value,
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
				Image:      "envcontainer/" + queries["1_project"].Value,
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

	err := os.MkdirAll(template.PathDefault, 0755)
	compose.check("error to create folders, check permissions", err)

	createFile := func(name string, data []byte) {
		compose.check("error to crete config files, check permissions", ioutil.WriteFile(name, data, 0644))
	}

	createFile(template.Dockerfile, []byte("FROM "+queries["2_image"].Value))
	createFile(template.Compose, data)

	envContent := fmt.Sprintf("COMPOSE_PROJECT_NAME=%s\nCOMPOSE_IGNORE_ORPHANS=True", queries["1_project"].Value)
	createFile(template.Env, []byte(envContent))

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
			panic(err)
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
