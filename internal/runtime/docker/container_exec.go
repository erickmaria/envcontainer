package docker

import (
	"context"
	"io"
	"os"

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

	docker.execInteractive(ctx, resp.ID)

	if options.AutoStop {
		return docker.Stop(ctx, options.ContainerName)
	}

	return nil
}

func (docker *Docker) execInteractive(ctx context.Context, containerID string) error {

	steam, err := docker.client.ContainerExecAttach(ctx, containerID, types.ExecStartCheck{
		Detach: false,
		Tty:    true,
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
