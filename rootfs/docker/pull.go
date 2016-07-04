package docker

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/distribution"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/reference"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

func (m *Manager) PullImage(ctx context.Context, image, tag string, metaHeaders map[string][]string, authConfig *types.AuthConfig, progressOutput progress.Output) error {
	ref, err := imageTagToReference(image, tag)
	if err != nil {
		return err
	}

	imagePullConfig := m.newImagePullConfig()
	imagePullConfig.MetaHeaders = metaHeaders
	imagePullConfig.AuthConfig = authConfig

	return m.pullWithProgress(ctx, ref, imagePullConfig, progressOutput)
}

func (m *Manager) newImagePullConfig() *distribution.ImagePullConfig {
	return &distribution.ImagePullConfig{
		DownloadManager:  m.downloadManager,
		ImageEventLogger: m.eventService.LogImageEvent,
		ImageStore:       m.imageStore,
		MetadataStore:    m.distributionMetadataStore,
		ReferenceStore:   m.referenceStore,
		RegistryService:  m.registryService,
	}
}

func (m *Manager) pullWithProgress(ctx context.Context, ref reference.Named, imagePullConfig *distribution.ImagePullConfig, progressOutput progress.Output) (err error) {
	progressChan := make(chan progress.Progress, 100)
	imagePullConfig.ProgressOutput = progress.ChanOutput(progressChan)
	ctx, cancelFunc := context.WithCancel(ctx)

	go func() {
		err = distribution.Pull(ctx, ref, imagePullConfig)
		close(progressChan)
	}()

	for prog := range progressChan {
		if progressOutput == nil {
			continue
		}

		if err := progressOutput.WriteProgress(prog); err != nil {
			logrus.Info("Canceling pull due to: %v", err)
			cancelFunc()
			progressOutput = nil
		}
	}

	return
}
