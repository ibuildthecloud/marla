package daemon

import (
	"errors"
	"io"
	"time"

	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/registry"
	"golang.org/x/net/context"
)

var ErrNotImplemented = errors.New("Not implemented")

func (d *Daemon) ContainerExecCreate(name string, config *types.ExecConfig) (string, error) {
	return "", ErrNotImplemented
}
func (d *Daemon) ContainerExecInspect(id string) (*backend.ExecInspect, error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) ContainerExecResize(name string, height, width int) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerExecStart(ctx context.Context, name string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer) error {
	return ErrNotImplemented
}
func (d *Daemon) ExecExists(name string) (bool, error) {
	return false, ErrNotImplemented
}
func (d *Daemon) ContainerArchivePath(name string, path string) (content io.ReadCloser, stat *types.ContainerPathStat, err error) {
	return nil, nil, ErrNotImplemented
}
func (d *Daemon) ContainerCopy(name string, res string) (io.ReadCloser, error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) ContainerExport(name string, out io.Writer) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerExtractToDir(name, path string, noOverwriteDirNonDir bool, content io.Reader) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerStatPath(name string, path string) (stat *types.ContainerPathStat, err error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) ContainerCreate(types.ContainerCreateConfig) (types.ContainerCreateResponse, error) {
	return types.ContainerCreateResponse{}, ErrNotImplemented
}
func (d *Daemon) ContainerKill(name string, sig uint64) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerPause(name string) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerRename(oldName, newName string) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerResize(name string, height, width int) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerRestart(name string, seconds int) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerRm(name string, config *types.ContainerRmConfig) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerStart(name string, hostConfig *container.HostConfig) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerStop(name string, seconds int) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerUnpause(name string) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerUpdate(name string, hostConfig *container.HostConfig) ([]string, error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) ContainerWait(name string, timeout time.Duration) (int, error) {
	return 0, ErrNotImplemented
}
func (d *Daemon) ContainerChanges(name string) ([]archive.Change, error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) ContainerInspect(name string, size bool, version string) (interface{}, error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) ContainerLogs(ctx context.Context, name string, config *backend.ContainerLogsConfig, started chan struct{}) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerStats(ctx context.Context, name string, config *backend.ContainerStatsConfig) error {
	return ErrNotImplemented
}
func (d *Daemon) ContainerTop(name string, psArgs string) (*types.ContainerProcessList, error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) Containers(config *types.ContainerListOptions) ([]*types.Container, error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) ContainerAttach(name string, c *backend.ContainerAttachConfig) error {
	return ErrNotImplemented
}

func (d *Daemon) Commit(name string, config *backend.ContainerCommitConfig) (imageID string, err error) {
	return "", ErrNotImplemented
}
func (d *Daemon) ImageDelete(imageRef string, force, prune bool) ([]types.ImageDelete, error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) ImageHistory(imageName string) ([]*types.ImageHistory, error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) Images(filterArgs string, filter string, all bool) ([]*types.Image, error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) LookupImage(name string) (*types.ImageInspect, error) {
	return nil, ErrNotImplemented
}
func (d *Daemon) TagImage(imageName, repository, tag string) error {
	return ErrNotImplemented
}
func (d *Daemon) LoadImage(inTar io.ReadCloser, outStream io.Writer, quiet bool) error {
	return ErrNotImplemented
}
func (d *Daemon) ImportImage(src string, repository, tag string, msg string, inConfig io.ReadCloser, outStream io.Writer, changes []string) error {
	return ErrNotImplemented
}
func (d *Daemon) ExportImage(names []string, outStream io.Writer) error {
	return ErrNotImplemented
}
func (d *Daemon) PushImage(ctx context.Context, image, tag string, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	return ErrNotImplemented
}
func (d *Daemon) SearchRegistryForImages(ctx context.Context, filtersArgs string, term string, limit int, authConfig *types.AuthConfig, metaHeaders map[string][]string) (*registry.SearchResults, error) {
	return nil, ErrNotImplemented
}
