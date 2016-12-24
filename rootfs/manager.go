package rootfs

import (
	"github.com/docker/docker/image"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

type RootCallback func(string) error

type Manager interface {
	PullImage(ctx context.Context, image, tag string, metaHeaders map[string][]string, authConfig *types.AuthConfig, progress progress.Output) error
	ResolveImageID(image string) (image.ID, error)
	//PrepareRoot(ctx context.Context, container model.Container, prepare RootCallback) error
}
