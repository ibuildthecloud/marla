package daemon

import (
	"github.com/docker/docker/image"
	"github.com/ibuildthecloud/marla/event"
	"github.com/ibuildthecloud/marla/rootfs"
	drootfs "github.com/ibuildthecloud/marla/rootfs/docker"
)

type Daemon struct {
	imageStore    *image.Store
	eventService  *event.Service
	rootFSManager rootfs.Manager
}

func New(config *Config) (*Daemon, error) {
	var err error

	d := &Daemon{
		eventService: event.New(),
	}

	if d.rootFSManager, err = drootfs.New(&config.DockerRootFS, d.eventService); err != nil {
		return nil, err
	}

	return d, nil
}
