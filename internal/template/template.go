package template

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/ErickMaria/envcontainer/internal/pkg/syscmd"
	"gopkg.in/yaml.v2"
)

var (
	paths = map[string]string{
		"home":        ".envcontainer",
		"cache":       ".envcontainer/cache",
		"tmp":         ".envcontainer/tmp",
		"dockerfiles": ".envcontainer/tmp/dockerfiles",
	}
)

type Envcontainer struct {
	Project struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Description string `yaml:"description"`
	} `yaml:"project"`
	Container struct {
		Ports []string `yaml:"ports"`
		Build string   `yaml:"build"`
	} `yaml:"container"`
}

func Initialization() error {

	_, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	err = syscmd.CreateDir(toSlice(paths))
	if err != nil {
		return err
	}

	return nil
}

func Unmarshal() (Envcontainer, error) {

	data, err := os.ReadFile(".envcontainer.yaml")
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
