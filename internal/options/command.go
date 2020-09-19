package options

import (
	"github.com/ErickMaria/envcontainer/internal/scan"
	"github.com/ErickMaria/envcontainer/internal/pkg"
)

const (
	PATH_DEFAULT = "../../.envcontainer/compose/env"
	DOCKERFILE = "../../.envcontainer/Dockerfile"
	DOCKER_COMPOSE = "../../.envcontainer/compose/docker-compose.yaml"
	ENV = "../../.envcontainer/compose/env/.env"
)

type Command struct {
} 

func (c *Command) Init(flags map[string]string) (string) {	
	pkg.CreateDir(PATH_DEFAULT)
	pkg.CreateFile([]string{DOCKERFILE, DOCKER_COMPOSE, ENV})

	dc := scan.DockerCompose{}
	dc.File(flags)

	return "init"
}

func (c *Command) Up() (string) {
	return "up"
}

func (c *Command) Build() (string) {
	return "build"
}

func (c *Command) Down() (string) {
	return "down"
}