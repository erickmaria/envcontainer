package options

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/config"
)

const (
	HOME           = ".envcontainer"
	PATH_DEFAULT   = ".envcontainer/compose"
	DOCKERFILE     = ".envcontainer/Dockerfile"
	DOCKER_COMPOSE = ".envcontainer/compose/docker-compose.yaml"
	ENV            = ".envcontainer/compose/.env"
)

const (
	INIT    string = "init"
	BUILD   string = "build"
	START   string = "start"
	STOP    string = "stop"
	DELETE  string = "delete"
	HELP    string = "help"
	VERSION string = "version"
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

	validate()

	values := flags.Values

	override := values["override"].valueBool

	if _, err := os.Stat(HOME); !os.IsNotExist(err) {

		if !*override {

			fmt.Print("already exists in this project, do you're have override? (yes/no): ")
			reader := bufio.NewReader(os.Stdin)
			confirmation, _, err := reader.ReadLine()
			check("error to read confirmation input, check input", err)

			v := string(confirmation)

			switch strings.ToLower(v) {
			case "yes":
				break
			case "no":
				return
			default:
				fmt.Println("values accepted are 'yes' or 'no'")
				return
			}

		}
	}

	var ports = []string{}

	if *values["listener"].valueString != "" {
		ports = append(ports, *values["listener"].valueString)
	}

	dc := config.DockerCompose{
		Version: "3.6",
		Services: config.Services{
			Environment: config.Environment{
				ContainerName: *values["project"].valueString,
				Volumes: []config.Volumes{
					config.Volumes{
						Type:   "bind",
						Source: "../../",
						Target: "/home/envcontainer/" + *values["project"].valueString,
					},
				},
				Build: config.Build{
					Dockerfile: "Dockerfile",
					Context:    "../",
				},
				Image:      "envcontainer/" + *values["project"].valueString,
				Ports:      ports,
				WorkingDir: "/home/envcontainer/" + *values["project"].valueString,
				EnvFile: []string{
					"$PWD/" + *values["envfile"].valueString,
				},
				Privileged: true,
				StdinOpen:  true,
				Tty:        true,
			},
		},
	}

	data := dc.Marshal()

	err := os.MkdirAll(PATH_DEFAULT, 0755)
	check("error to create folders, check permissions", err)

	createFile := func(name string, data []byte) {
		check("error to crete config files, check permissions", ioutil.WriteFile(name, data, 0644))
	}

	createFile(DOCKERFILE, []byte("FROM "+*values["image"].valueString))
	createFile(DOCKER_COMPOSE, data)

	envContent := fmt.Sprintf("COMPOSE_PROJECT_NAME=%s\nCOMPOSE_IGNORE_ORPHANS=True", *values["project"].valueString)
	createFile(ENV, []byte(envContent))

	if !*values["no-build"].valueBool {
		Build()
	}

	fmt.Println("envcontainer initialized!")
}

func Build() {

	validate()

	command(
		"docker-compose",
		"-f",
		DOCKER_COMPOSE,
		"build",
	)

	fmt.Println("envcontainer: build successful!")

}

func Start(flags Flag) {

	validate()

	dc := config.DockerComposeConfig.Unmarshal(DOCKER_COMPOSE)

	command(
		"docker-compose",
		"-f",
		DOCKER_COMPOSE,
		"up",
		"-d",
		"--no-log-prefix",
	)

	// command("clear")

	command(
		"docker",
		"exec",
		"-it",
		dc.Services.Environment.ContainerName,
		*flags.Values["shell"].valueString,
	)
}

func Stop() {

	validate()

	command(
		"docker-compose",
		"-f",
		DOCKER_COMPOSE,
		"down",
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

	if !*autoApprove.valueBool {
		fmt.Print("do you're have sure? (yes/no): ")
		reader := bufio.NewReader(os.Stdin)
		confirmation, _, err := reader.ReadLine()
		check("error to read confirmation input, check input", err)

		v := string(confirmation)

		switch strings.ToLower(v) {
		case "yes":
			break
		case "no":
			return
		default:
			fmt.Println("envcontainer: values accepted are 'yes' or 'no'")
			return
		}
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
	check("cannot remove files, check permissions", err)

	fmt.Println("envcontainer: configuration deleted")

}

func validate() {
	stdout, err := exec.Command("docker", "info").Output()

	if strings.Contains(string(stdout), "unix:///var/run/docker.sock") {

		check("Is the docker daemon running?", err)
	}
}

func command(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		check("command failed, check envcontainer configs.", err)
	}
}

func check(message string, e error) {
	if e != nil {
		fmt.Println("envcontainer: " + message + "\n" + e.Error())
		os.Exit(0)
	}
}
