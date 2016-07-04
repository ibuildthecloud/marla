package rootfs

import (
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

type Manager interface {
	PullImage(ctx context.Context, image, tag string, metaHeaders map[string][]string, authConfig *types.AuthConfig, progress progress.Output) error
}