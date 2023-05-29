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
		"cache":       "/tmp/.envcontainer/cache",
		"tmp":         "/tmp/.envcontainer/tmp",
		"dockerfiles": "/tmp/.envcontainer/tmp/dockerfiles",
	}
	fileLocation string = ".envcontainer.yaml"
)

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
	AlwaysUpdate bool     `yaml:"always-update"`
	AutoStop     bool     `yaml:"auto-stop"`
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

	envcontainer.Project.Name = strings.ToLower(envcontainer.Project.Name)
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
		dockerfile := paths["dockerfiles"] + "/" + "Dockerfile." + envcontainer.Project.Name + "-" + envcontainer.Project.Version
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

func toSlice(maps map[string]string) []string {

	values := []string{}
	for _, v := range maps {
		values = append(values, v)
	}

	return values
}

func GetCacheDir() string {
	return paths["cache"]
}
