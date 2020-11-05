package options

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
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
	BUILD  string = "build"
	START  string = "start"
	STOP   string = "stop"
	DELETE string = "delete"
	HELP   string = "help"
)

type CommandConfig map[string]Command

type Command struct {
	Flags Flag
	Exec  func()
	Desc  string
}

func NewCommand(cc CommandConfig) (*Command, CommandConfig) {

	if isEmpty() {
		Help(cc)
		os.Exit(0)
	}

	if !contains(cc) {
		fmt.Printf("envcontainer: '%s' is not a envcontainer command\n%s\n", os.Args[1], cc[HELP].Desc)
		os.Exit(0)
	}

	ccAux := cc[os.Args[1]]
	return &ccAux, cc
}

func (c Command) Listener() {
	c.Flags.Register()
	c.Exec()
}

func isEmpty() bool {
	if len(os.Args) < 2 {
		return true
	}
	return false
}

func contains(cc CommandConfig) bool {

	for k := range cc {
		if strings.ToLower(os.Args[1]) == k {
			return true
		}
	}

	return false
}

func Init(flags Flag) {

	values := flags.Values

	override := values["override"]

	// fmt.Println(*override.value)

	if _, err := os.Stat(HOME); !os.IsNotExist(err) {

		if *override.value != "yes" {

			fmt.Print("envcontainer has exists in this project, do you're have override? (yes/no): ")
			reader := bufio.NewReader(os.Stdin)
			confirmation, _, err := reader.ReadLine()
			check("envcontainer: error to read confirmation input, check input", err)

			v := string(confirmation)
			override.value = &v

			switch strings.ToLower(*override.value) {
			case "yes":
				break
			case "no":
				return
			default:
				fmt.Println("envcontainer: values accepted are 'yes' or 'no'")
				return
			}
		}
	}

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
	check("envcontainer: error to create folders, check permissions", err)

	createFile := func(name string, data []byte) {
		check("envcontainer: error to crete config files, check permissions", ioutil.WriteFile(name, data, 0644))
	}

	createFile(DOCKERFILE, []byte("FROM "+*values["image"].value))
	createFile(DOCKER_COMPOSE, data)
	createFile(ENV, []byte(""))

	noBuild, err2 := strconv.ParseBool(*values["no-build"].value)
	if err2 != nil {
		fmt.Println("envcontainer: values accepted are 'true' or 'false'")
		os.Exit(0)
	}

	if !noBuild {
		Build()
	}

	fmt.Println("envcontainer initialized!")
}

func Build() {

	command(
		"docker-compose",
		"-f",
		DOCKER_COMPOSE,
		"up",
		"-d",
		"--build",
	)

	fmt.Println("envcontainer: build successful!")

}

func Start(flags Flag) {

	dc := config.DockerComposeConfig.Unmarshal(DOCKER_COMPOSE)

	command(
		"docker",
		"start",
		dc.Services.Environment.ContainerName,
	)

	command(
		"docker",
		"exec",
		dc.Services.Environment.ContainerName,
		"echo",
		"envcontainer: connected!",
	)

	command(
		"docker",
		"exec",
		"-it",
		dc.Services.Environment.ContainerName,
		*flags.Values["shell"].value,
	)
}

func command(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stdout = os.Stderr
	err := cmd.Run()

	if err != nil {
		check("envcontainer: command failed, check envcontainer configs and permisstions", err)
	}
}

func Stop() {

	command(
		"docker-compose",
		"-f",
		DOCKER_COMPOSE,
		"kill",
	)

	fmt.Println("envcontainer: stopped")

}

func Help(descs map[string]Command) {

	fmt.Println("\nUsage: envcontainer COMMAND --FLAGS")

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
		check("envcontainer: error to read confirmation input, check input", err)

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

	command(
		"docker-compose",
		"-f",
		DOCKER_COMPOSE,
		"down",
		"--rmi",
		"all",
	)

	err := os.RemoveAll(".envcontainer")
	check("envcontainer: cannot remove files, check permissions", err)

	fmt.Println("envcontainer: configuration deleted")

}

func check(message string, e error) {
	if e != nil {
		panic(message + ": " + e.Error())
	}
}
