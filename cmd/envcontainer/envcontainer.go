package main

import (
	"context"

	"github.com/ErickMaria/envcontainer/internal/runtime/docker"
	"github.com/ErickMaria/envcontainer/pkg/cli"
)

var cmd *cli.Command
var cmds cli.CommandConfig
var path string

// # DOCKER API
var ctx = context.Background()
var container = docker.NewDocker()

func main() {
	Root().Listener()
}
