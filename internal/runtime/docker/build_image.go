package docker

import (
	"context"
	"io"
	"os"
	"strings"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/archive"
)

func (docker *Docker) Build(ctx context.Context, options runtimeTypes.BuildOptions) error {

	buildCtx, err := archive.TarWithOptions(options.BuildContext, &archive.TarOptions{})
	if err != nil {
		return err
	}

	imageBuildResponse, err := docker.client.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Tags:       []string{"envcontainer/" + options.ImageName},
		Dockerfile: "Dockerfile",
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

func (docker *Docker) pullImage(ctx context.Context, image string) error {

	out, err := docker.client.ImagePull(ctx, image, types.ImagePullOptions{})
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

func (docker *Docker) checkIfImageExists(ctx context.Context, image string) (bool, error) {

	if len(strings.Split(image, ":")) != 2 {
		image = image + ":latest"
	}

	images, err := docker.client.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: image,
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
