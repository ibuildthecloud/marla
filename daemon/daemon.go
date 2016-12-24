package daemon

import (
	"github.com/ibuildthecloud/marla/container"
	"github.com/ibuildthecloud/marla/event"
	"github.com/ibuildthecloud/marla/rootfs"
	drootfs "github.com/ibuildthecloud/marla/rootfs/docker"
)

type Daemon struct {
	//imageStore    *image.Store
	containerStore *container.Store
	eventService   *event.Service
	rootFSManager  rootfs.Manager
}

func New(config *Config) (*Daemon, error) {
	var err error

	d := &Daemon{
		eventService: event.New(),
	}

	if d.containerStore, err = container.NewStore(&config.Container); err != nil {
		return nil, err
	}

	if d.rootFSManager, err = drootfs.New(&config.DockerRootFS, d.eventService); err != nil {
		return nil, err
	}

	return d, nil
}
