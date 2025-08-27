package docker

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	pkgTypes "github.com/ErickMaria/envcontainer/internal/pkg/types"
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

	err := docker.Down(ctx, runtimeTypes.ContainerOptions{
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

func (docker *Docker) getContainer(ctx context.Context, labels map[string]string) (types.Container, error) {

	kvLabels := []filters.KeyValuePair{}

	for k, v := range labels {
		if strings.Contains(k, "project-name") || strings.Contains(k, "project-path") {
			kvLabels = append(kvLabels, filters.KeyValuePair{
				Key:   "label",
				Value: k + "=" + v,
			})
		}
	}

	containers, err := docker.client.ContainerList(ctx, container.ListOptions{
		Limit:   1,
		Filters: filters.NewArgs(kvLabels...),
	})

	// containers, err := docker.client.ContainerList(ctx, container.ListOptions{
	// 	Limit: 1,
	// 	Filters: filters.NewArgs(filters.KeyValuePair{
	// 		Key:   "name",
	// 		Value: containerName,
	// 	}),
	// })

	if len(containers) == 0 {
		return types.Container{}, nil
	}

	if err != nil {
		return types.Container{}, err
	}

	return containers[0], nil
}

func (docker *Docker) getNetwork(ctx context.Context, labels map[string]string) ([]types.NetworkResource, error) {

	kvLabels := []filters.KeyValuePair{}

	for k, v := range labels {
		if strings.Contains(k, "project-name") || strings.Contains(k, "project-path") {
			kvLabels = append(kvLabels, filters.KeyValuePair{
				Key:   "label",
				Value: k + "=" + v,
			})
		}
	}

	networks, err := docker.client.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(kvLabels...),
	})

	if len(networks) == 0 {
		return []types.NetworkResource{}, nil
	}

	if err != nil {
		return []types.NetworkResource{}, err
	}

	return networks, nil
}

// func (docker *Docker) buildMount(defaultMountDir string, mountStr []string) []mount.Mount {

// 	mounts := []mount.Mount{}
// 	mountPrefix := "envcontainer"

// 	for k, v := range mountStr {

// 		mountSplited := strings.Split(v, ":")

// 		// fmt.Println(len(mountSplited), mountSplited)

// 		switch len(mountSplited) {
// 		case 1:

// 			target := mountSplited[1-len(mountSplited)]
// 			mounts = append(mounts, mount.Mount{
// 				Type:   mount.TypeVolume,
// 				Source: mountPrefix + "-" + strings.Trim(target, "/"),
// 				Target: target,
// 			})
// 		case 2:

// 			mountType := strings.ToLower(mountSplited[1])

// 			if strings.Contains(mountType, string(mount.TypeBind)) {

// 				source := removeDuplicateSlashes(defaultMountDir + mountSplited[0])

// 				fmt.Println(source)

// 				mounts = append(mounts, mount.Mount{
// 					Type:   mount.TypeBind,
// 					Source: source,
// 					Target: mountSplited[0],
// 					BindOptions: &mount.BindOptions{
// 						CreateMountpoint: true,
// 					},
// 				})

// 			}

// 			// when type volume or not define
// 			mounts = append(mounts, mount.Mount{
// 				Type:   mount.TypeVolume,
// 				Source: mountPrefix + "-" + mountSplited[0],
// 				Target: mountSplited[1],
// 			})

// 		case 3:
// 			mountType := strings.ToLower(mountSplited[2])

// 			if strings.Contains(mountType, string(mount.TypeVolume)) {

// 				mounts = append(mounts, mount.Mount{
// 					Type:   mount.TypeVolume,
// 					Source: mountPrefix + "-" + mountSplited[0],
// 					Target: mountSplited[1],
// 				})
// 				continue
// 			} else if strings.Contains(mountType, string(mount.TypeBind)) {

// 				mounts = append(mounts, mount.Mount{
// 					Type:   mount.TypeBind,
// 					Source: mountSplited[0],
// 					Target: mountSplited[1],
// 				})
// 				continue
// 			}

// 			mountPatternNotMatchError(mountStr[k])

// 		default:
// 			mountPatternNotMatchError(mountStr[k])
// 		}

// 	}

// 	return mounts
// }

func (docker *Docker) buildMount(defaultMountDir string, mounts []pkgTypes.Mount, labels map[string]string) []mount.Mount {

	dMounts := []mount.Mount{}

	for _, v := range mounts {

		if v.Type == "volume" {
			dMounts = append(dMounts, mount.Mount{
				Type:     mount.Type(v.Type),
				Source:   v.Source,
				Target:   v.Target,
				ReadOnly: v.Readonly,
				VolumeOptions: &mount.VolumeOptions{
					DriverConfig: &mount.Driver{},
					Labels:       labels,
				},
			})
			continue
		}

		dMounts = append(dMounts, mount.Mount{
			Type:     mount.Type(v.Type),
			Source:   v.Source,
			Target:   v.Target,
			ReadOnly: v.Readonly,
		})

	}

	return dMounts
}

// func removeDuplicateSlashes(s string) string {
// 	re := regexp.MustCompile(`/{2,}`)
// 	return re.ReplaceAllString(s, "/")
// }

// func mountPatternNotMatchError(mount string) {
// 	panic("envcontainer: mount " + mount + " does not match the pattern.\n")
// }

func (docker *Docker) code(ctx context.Context, containerID string, port string, options runtimeTypes.ContainerOptions) error {

	inspect, err := docker.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return err
	}

	if inspect.State.Status == "exited" {
		docker.Down(ctx, runtimeTypes.ContainerOptions{
			ContainerName: strings.Split(options.ContainerName, "-")[0],
			HostDirToBind: options.HostDirToBind,
		})
	}

	address := inspect.NetworkSettings.IPAddress
	if strings.ToLower(options.NetworkMode) == "host" {
		address = "0.0.0.0"
	}

	var host = fmt.Sprintf("%s:%s", address, port)
	for i, p := range inspect.NetworkSettings.Ports {
		if strings.Contains("22", i.Port()) {
			host = fmt.Sprintf("%s:%s", address, p[0].HostPort)
		}
	}

	if !docker.isPortAvailable(address, port, 3*time.Second) {
		fmt.Println("port " + port + " is not available. Try to use --port flag or try again")
		os.Exit(1)
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

func (docker *Docker) isPortAvailable(host string, port string, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%s", host, port)

	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}
