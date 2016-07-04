package daemon

import (
	"io"

	"github.com/docker/docker/pkg/streamformatter"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

func (d *Daemon) PullImage(ctx context.Context, image, tag string, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	progressOutput := streamformatter.NewJSONStreamFormatter().NewProgressOutput(outStream, false)
	return d.rootFSManager.PullImage(ctx, image, tag, metaHeaders, authConfig, progressOutput)
}
