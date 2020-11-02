package options

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/config"
)

const (
	HOME           = ".envcontainer"
	PATH_DEFAULT   = ".envcontainer/compose/env"
	DOCKERFILE     = ".envcontainer/Dockerfile"
	DOCKER_COMPOSE = ".envcontainer/compose/docker-compose.yaml"
	ENV            = ".envcontainer/compose/env/.variables"
)

const (
	INIT   string = "init"
	RUN    string = "run"
	STOP   string = "stop"
	DELETE string = "delete"
	HELP   string = "help"
)

type Command struct {
	Flags Flag
	Exec  func()
	Desc  string
}

func Init(flags Flag) {

	values := flags.Values

	var ports = []string{}

	if *values["listener"].value != "" {
		ports = append(ports, *values["listener"].value)
	}

	dc := config.DockerCompose{
		Version: "3.6",
		Services: config.Services{
			Environment: config.Environment{
				ContainerName: *values["project"].value,
				Volumes: []config.Volumes{
					config.Volumes{
						Type:   "bind",
						Source: "../../",
						Target: "/home/envcontainer/" + *values["project"].value,
					},
				},
				Build: config.Build{
					Dockerfile: "Dockerfile",
					Context:    "../",
				},
				Ports:      ports,
				WorkingDir: "/home/envcontainer/" + *values["project"].value,
				EnvFile: []string{
					*values["envfile"].value,
				},
				Privileged: true,
				StdinOpen:  true,
				Tty:        true,
			},
		},
	}

	data := dc.Marshal()

	err := os.MkdirAll(PATH_DEFAULT, 0755)
	check(err)

	createFile := func(name string, data []byte) {
		check(ioutil.WriteFile(name, data, 0644))
	}

	createFile(DOCKERFILE, []byte("FROM "+*values["image"].value))
	createFile(DOCKER_COMPOSE, data)
	createFile(ENV, []byte(""))
}

func Help(descs map[string]Command) {

	fmt.Println("\nUsage:  envcontainer COMMAND --FLAGS")

	fmt.Println("\nCommands")

	for commandKey, comandValue := range descs {
		fmt.Printf("%s:     \t%v\n", commandKey, descs[commandKey].Desc)
		for flagKey, flagValue := range comandValue.Flags.Values {
			fmt.Printf("    --%s:     \t%v\n", flagKey, flagValue.Description)
		}

	}

	fmt.Println()
}

func Delete(flags Flag) {

	values := flags.Values
	autoApprove := values["auto-approve"]

	if *autoApprove.value == "" {
		fmt.Print("do you're have sure? (yes/no): ")
		reader := bufio.NewReader(os.Stdin)
		confirmation, _, err := reader.ReadLine()
		check(err)

		v := string(confirmation)
		autoApprove.value = &v
	}

	switch strings.ToLower(*autoApprove.value) {
	case "yes":
		break
	case "no":
		return
	default:
		fmt.Println("envcontainer: values accepted are 'yes' or 'no'")
		return
	}
	os.RemoveAll(".envcontainer")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
