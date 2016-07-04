package daemon

import "github.com/ibuildthecloud/marla/rootfs/docker"

type Config struct {
	Root         string
	DockerRootFS docker.Config
}
