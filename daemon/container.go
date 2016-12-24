package daemon

import "github.com/docker/engine-api/types"

func (d *Daemon) ContainerCreate(params types.ContainerCreateConfig) (types.ContainerCreateResponse, error) {
	result := types.ContainerCreateResponse{}

	imageID, err := d.rootFSManager.ResolveImageID(params.Config.Image)
	if err != nil {
		return result, err
	}
	id, err := d.containerStore.Register(params, imageID.String())
	result.ID = id
	return result, err
}
