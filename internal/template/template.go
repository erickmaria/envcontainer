package template

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/pkg/syscmd"
	"gopkg.in/yaml.v2"
)

var (
	paths = map[string]string{
		"home":        "/tmp/.envcontainer",
		"dockerfiles": "/tmp/.envcontainer/dockerfiles",
	}
	fileLocation string = ".envcontainer.yaml"
)

type Type string

// Type constants
const (
	TypeBind   Type = "bind"
	TypeVolume Type = "volume"
)

type EnvcontainerMounts struct {
	Type   Type   `yaml:"type"`
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type Envcontainer struct {
	Project struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Description string `yaml:"description"`
	} `yaml:"project"`
	Container struct {
		User  string   `yaml:"user"`
		Ports []string `yaml:"ports"`
		Build string   `yaml:"build"`
	} `yaml:"container"`
	AlwaysUpdate bool `yaml:"always-update"`
	AutoStop     bool `yaml:"auto-stop"`
	mountDir     string
	Mounts       []string `yaml:"mounts"`
}

func Initialization() error {

	err := syscmd.CreateDir(toSlice(paths))
	if err != nil {
		return err
	}

	return nil
}

func Unmarshal() (Envcontainer, error) {

	data, err := os.ReadFile(fileLocation)
	if err != nil {
		return Envcontainer{}, err
	}

	var envcontainer Envcontainer
	err = yaml.Unmarshal(data, &envcontainer)
	if err != nil {
		return Envcontainer{}, err
	}

	envcontainer.Project.Name = strings.ReplaceAll(strings.ToLower(envcontainer.Project.Name), " ", "-")
	envcontainer.Container.Build, err = tmpDockerfile(envcontainer)
	if err != nil {
		return Envcontainer{}, err
	}

	err = validate(envcontainer)
	if err != nil {
		return Envcontainer{}, err
	}

	return envcontainer, nil
}

func UnmarshalWithFile(location string) (Envcontainer, error) {

	fileLocation = location

	return Unmarshal()
}

func validate(envcontainer Envcontainer) error {

	if len(envcontainer.Container.Ports) > 0 {

		regex := regexp.MustCompile("^(\\d+)(:?)(\\d+)$")

		for _, v := range envcontainer.Container.Ports {

			if ok := regex.MatchString(v); !ok {
				return errors.New("port " + v + " is invalid")
			}
		}
	}

	return nil
}

func tmpDockerfile(envcontainer Envcontainer) (string, error) {

	_, err := os.ReadFile(envcontainer.Container.Build)
	if err != nil {
		dockerfilePath := envcontainer.GetTmpDockerfileDir()

		err = syscmd.CreateDir([]string{dockerfilePath})
		if err != nil {
			return "", err
		}
		dockerfile := dockerfilePath + "/Dockerfile"
		err = syscmd.CreateFile(dockerfile, []byte(envcontainer.Container.Build))
		if err != nil {
			return "", err
		}

		user := envcontainer.Container.User
		if user != "" {
			syscmd.AppendFile(dockerfile, []byte(`
RUN apt-get update && apt-get install sudo
RUN useradd -ms /bin/bash `+user+` && echo "`+user+` ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
USER `+user+`
		`))
		}

		return dockerfile, nil

	}
	return envcontainer.Container.Build, nil
}

func (envcontainer Envcontainer) GetTmpDockerfileDir() string {
	return paths["dockerfiles"] + "/" + envcontainer.Project.Name + "/" + envcontainer.Project.Version
}

func toSlice(maps map[string]string) []string {

	values := []string{}
	for _, v := range maps {
		values = append(values, v)
	}

	return values
}

func sliceDeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func (envcontainer *Envcontainer) SetMountDir(mountDir string) {

	envcontainer.mountDir = mountDir

}

func (envcontainer *Envcontainer) GetMountDir() string {
	return envcontainer.mountDir
}
