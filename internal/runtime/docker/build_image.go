package docker

import (
	"context"
	"io"
	"os"
	"strings"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/archive"
)

func (docker *Docker) Build(ctx context.Context, options runtimeTypes.BuildOptions) error {

	buildCtx, err := archive.TarWithOptions(options.BuildContext, &archive.TarOptions{})
	if err != nil {
		return err
	}

	imageBuildResponse, err := docker.client.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Tags:        []string{"envcontainer/" + options.ImageName},
		NetworkMode: options.NetworkMode,
		Dockerfile:  "Dockerfile",
	})
	if err != nil {
		return err
	}
	defer imageBuildResponse.Body.Close()

	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return err
	}

	return nil
}

func (docker *Docker) pullImage(ctx context.Context, imageName string) error {

	out, err := docker.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return err
	}

	return nil
}

func (docker *Docker) checkIfImageExists(ctx context.Context, imageName string) (bool, error) {

	if len(strings.Split(imageName, ":")) != 2 {
		imageName = imageName + ":latest"
	}

	images, err := docker.client.ImageList(ctx, image.ListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: imageName,
		}),
	})
	if err != nil {
		return false, err
	}

	if len(images) == 0 {
		return false, nil
	}

	return true, nil
}
