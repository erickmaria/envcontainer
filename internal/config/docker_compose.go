package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type DockerCompose struct {
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

var DockerComposeConfig *DockerCompose

func (dc *DockerCompose) Marshal() []byte {

	data, err := yaml.Marshal(dc)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return data
}

func (dc *DockerCompose) Unmarshal(path string) DockerCompose {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	optDc := DockerCompose{}

	err2 := yaml.Unmarshal(data, &optDc)

	if err2 != nil {
		log.Fatalf("error: %v", err)
	}

	return optDc
}
