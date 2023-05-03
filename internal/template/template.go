package template

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
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
