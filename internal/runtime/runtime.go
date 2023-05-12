package runtime

import (
	"context"

	"github.com/ErickMaria/envcontainer/internal/runtime/types"
)

type ContainerRuntime interface {
	Start(context.Context, types.ContainerOptions) error
	AlwaysUpdate(context.Context, types.BuildOptions) error
}
