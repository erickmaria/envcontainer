package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

var shells = []string{"/bin/bash", "/bin/sh"}

type Docker struct {
	client *client.Client
}

func NewDocker() *Docker {

	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &Docker{
		client: client,
	}
}

func (docker *Docker) AlwaysUpdate(ctx context.Context, options runtimeTypes.BuildOptions) error {

	err := docker.Stop(ctx, runtimeTypes.ContainerOptions{
		ContainerName: options.ImageName,
	})
	if err != nil {
		panic(err)
	}

	return docker.Build(ctx, options)

}

func (docker *Docker) addContainerSuffix(options *runtimeTypes.ContainerOptions) {

	if !options.NoContainerSuffix {
		pathSplit := strings.Split(options.HostDirToBind, "/")
		containerNameSuffix := pathSplit[len(pathSplit)-1]
		options.ContainerName = strings.ToLower(options.ContainerName + "-" + containerNameSuffix)
	}
}

func (docker *Docker) getContainer(ctx context.Context, containerName string) (types.Container, error) {

	containers, err := docker.client.ContainerList(ctx, container.ListOptions{
		Limit: 1,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: containerName,
		}),
	})

	if len(containers) == 0 {
		return types.Container{}, nil
	}

	if err != nil {
		return types.Container{}, err
	}

	return containers[0], nil
}

func (docker *Docker) buildMount(defaultMountDir string, mountStr []string) []mount.Mount {

	mounts := []mount.Mount{}
	mountPrefix := "envcontainer"

	for k, v := range mountStr {

		mountSplited := strings.Split(v, ":")

		// fmt.Println(len(mountSplited), mountSplited)

		switch len(mountSplited) {
		case 1:

			target := mountSplited[1-len(mountSplited)]
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeVolume,
				Source: mountPrefix + "-" + strings.Trim(target, "/"),
				Target: target,
			})
		case 2:

			mountType := strings.ToLower(mountSplited[1])

			if strings.Contains(mountType, string(mount.TypeBind)) {

				source := removeDuplicateSlashes(defaultMountDir + mountSplited[0])

				fmt.Println(source)

				mounts = append(mounts, mount.Mount{
					Type:   mount.TypeBind,
					Source: source,
					Target: mountSplited[0],
					BindOptions: &mount.BindOptions{
						CreateMountpoint: true,
					},
				})

			}

			// when type volume or not define
			mounts = append(mounts, mount.Mount{
				Type:   mount.TypeVolume,
				Source: mountPrefix + "-" + mountSplited[0],
				Target: mountSplited[1],
			})

		case 3:
			mountType := strings.ToLower(mountSplited[2])

			if strings.Contains(mountType, string(mount.TypeVolume)) {

				mounts = append(mounts, mount.Mount{
					Type:   mount.TypeVolume,
					Source: mountPrefix + "-" + mountSplited[0],
					Target: mountSplited[1],
				})
				continue
			} else if strings.Contains(mountType, string(mount.TypeBind)) {

				mounts = append(mounts, mount.Mount{
					Type:   mount.TypeBind,
					Source: mountSplited[0],
					Target: mountSplited[1],
				})
				continue
			}

			mountPatternNotMatchError(mountStr[k])

		default:
			mountPatternNotMatchError(mountStr[k])
		}

	}

	return mounts
}

func removeDuplicateSlashes(s string) string {
	re := regexp.MustCompile(`/{2,}`)
	return re.ReplaceAllString(s, "/")
}

func mountPatternNotMatchError(mount string) {
	panic("envcontainer: mount " + mount + " does not match the pattern.\n")
}

func (docker *Docker) code(ctx context.Context, containerID string, options runtimeTypes.ContainerOptions) error {

	inspect, err := docker.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return err
	}

	if inspect.State.Status == "exited" {
		docker.Stop(ctx, runtimeTypes.ContainerOptions{
			ContainerName: strings.Split(options.ContainerName, "-")[0],
			HostDirToBind: options.HostDirToBind,
		})
	}

	var host = inspect.NetworkSettings.IPAddress + ":22"

	for i, p := range inspect.NetworkSettings.Ports {
		if strings.Contains("22", i.Port()) {
			host = "localhost:" + p[0].HostPort
		}
	}

	connection := fmt.Sprint("vscode-remote://ssh-remote+"+inspect.Config.User+"@"+host, "/home/"+options.ContainerName)

	fmt.Println("conecting to container " + containerID + " from the remote " + connection)

	cmd := exec.Command("code", "--folder-uri", connection)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	return nil
}
