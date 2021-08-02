package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

type Compose struct {
	Version  string   `yaml:"version"`
	Services Services `yaml:"services"`
}
type Build struct {
	Dockerfile string `yaml:"dockerfile"`
	Context    string `yaml:"context"`
}
type Volumes struct {
	Type   string `yaml:"type"`
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}
type Environment struct {
	ContainerName string    `yaml:"container_name"`
	Build         Build     `yaml:"build"`
	Image         string    `yaml:"image"`
	Ports         []string  `yaml:"ports"`
	WorkingDir    string    `yaml:"working_dir"`
	EnvFile       []string  `yaml:"env_file"`
	Volumes       []Volumes `yaml:"volumes"`
	Privileged    bool      `yaml:"privileged"`
	StdinOpen     bool      `yaml:"stdin_open"`
	Tty           bool      `yaml:"tty"`
}
type Services struct {
	Environment Environment `yaml:"environment"`
}

var ComposeConfig *Compose

func (compose *Compose) Marshal() []byte {

	data, err := yaml.Marshal(compose)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return data
}

func (compose *Compose) Unmarshal(path string) Compose {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	optCompose := Compose{}

	err2 := yaml.Unmarshal(data, &optCompose)

	if err2 != nil {
		log.Fatalf("error: %v", err)
	}

	return optCompose
}

func (compose *Compose) Build() {

	compose.validate()

	compose.command(
		"docker-compose",
		"-f",
		COMPOSE,
		"build",
	)

	fmt.Println("envcontainer: build successful!")

}

func (compose *Compose) Up(shell string) {

	compose.validate()

	dc := ComposeConfig.Unmarshal(COMPOSE)

	compose.command(
		"docker-compose",
		"-f",
		COMPOSE,
		"up",
		"-d",
	)

	compose.command(
		"docker",
		"exec",
		"-it",
		dc.Services.Environment.ContainerName,
		shell,
	)
}

func (compose *Compose) Down() {

	compose.validate()

	compose.command(
		"docker-compose",
		"-f",
		COMPOSE,
		"down",
	)

	fmt.Println("envcontainer: stopped")

}

func (compose *Compose) Delete() {

	compose.command(
		"docker-compose",
		"-f",
		COMPOSE,
		"down",
		"--rmi",
		"all",
	)

	err := os.RemoveAll(".envcontainer")
	compose.check("cannot remove files, check permissions", err)

	fmt.Println("envcontainer: configuration deleted")

}

func (compose *Compose) command(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		compose.check("command failed, check envcontainer configs.", err)
	}
}

func (compose *Compose) validate() {
	stdout, err := exec.Command("docker", "info").Output()

	if strings.Contains(string(stdout), "unix:///var/run/docker.sock") {

		compose.check("Is the docker daemon running?", err)
	}
}

func (compose *Compose) check(message string, e error) {
	if e != nil {
		fmt.Println("envcontainer: " + message + "\n" + e.Error())
		os.Exit(0)
	}
}
