package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	runtimeTypes "github.com/ErickMaria/envcontainer/internal/runtime/types"
	"github.com/docker/docker/api/types"
	"github.com/moby/term"
)

func (docker *Docker) exec(ctx context.Context, containerID string, options runtimeTypes.ContainerOptions) error {

	resp, err := docker.client.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Privileged:   true,
		Tty:          true,
		Cmd:          options.Commands,
	})
	if err != nil {
		return err
	}

	if err := docker.execInteractive(ctx, resp.ID); err != nil {
		docker.Stop(ctx, runtimeTypes.ContainerOptions{
			ContainerName: strings.Split(options.ContainerName, "-")[0] ,
			HostDirToBind: options.HostDirToBind,
		})
		return err
	}

	if options.AutoStop {
		return docker.Stop(ctx, runtimeTypes.ContainerOptions{
			ContainerName: strings.Split(options.ContainerName, "-")[0] ,
			HostDirToBind: options.HostDirToBind,
		})
	}

	return nil
}

func (docker *Docker) execInteractive(ctx context.Context, containerID string) error {

	height, width, err := getTerminalSize()
	if err != nil {
		fmt.Println("Error getting terminal size: ", err)
	}

	steam, err := docker.client.ContainerExecAttach(ctx, containerID, types.ExecStartCheck{
		ConsoleSize: &[2]uint{height, width},
		Detach:      false,
		Tty:         true,
	})

	if err != nil {
		return err
	}

	defer steam.Close()

	state, err := term.MakeRaw(os.Stdin.Fd())
	if err != nil {
		panic(err)
	}
	defer term.RestoreTerminal(os.Stdin.Fd(), state)

	go func() {
		_, err = io.Copy(steam.Conn, os.Stdin)
		if err != nil {
			panic(err)
		}
	}()

	_, err = io.Copy(os.Stdout, steam.Reader)
	if err != nil {
		return err
	}

	return nil
}

func getTerminalSize() (uint, uint, error) {
	ws, err := term.GetWinsize(os.Stderr.Fd())
	if err != nil {
		if ws == nil {
			return 0, 0, err
		}
	}
	return uint(ws.Height), uint(ws.Width), nil
}
