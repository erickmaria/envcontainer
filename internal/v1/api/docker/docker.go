package docker

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type Docker struct {
	cli *client.Client
}

func NewDocker() *Docker {
	
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &Docker{
		cli: cli,
	}
}

func (docker *Docker) Build(ctx context.Context) error {

	buildCtx, _ := archive.TarWithOptions(".envcontainer", &archive.TarOptions{})

	imageBuildResponse, err := docker.cli.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Tags: []string{"envcontainer/envcontainer"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		panic(err)
	}

	return nil
}
