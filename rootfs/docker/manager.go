package docker

import (
	"path/filepath"

	"github.com/docker/docker/distribution/metadata"
	"github.com/docker/docker/distribution/xfer"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/reference"
	"github.com/docker/docker/registry"
	"github.com/ibuildthecloud/marla/event"
	"github.com/pkg/errors"
)

type Manager struct {
	layerStore                layer.Store
	imageStore                image.Store
	distributionMetadataStore metadata.Store
	downloadManager           *xfer.LayerDownloadManager
	uploadManager             *xfer.LayerUploadManager
	referenceStore            reference.Store
	registryService           registry.Service
	eventService              *event.Service
}

type Config struct {
	Root      string
	Graph     GraphConfig
	Transfer  TransferConfig
	RemapRoot UserConfig
	Registry  registry.ServiceOptions
}

type GraphConfig struct {
	Driver  string
	Options []string
}

type UserConfig struct {
	Username string
	Group    string
}

type TransferConfig struct {
	MaxConcurrentUploads   int
	MaxConcurrentDownloads int
}

func New(config *Config, es *event.Service) (*Manager, error) {
	var (
		err error
		m   Manager
	)

	m.eventService = es

	if err = m.setupLayerStore(config); err != nil {
		return nil, errors.Wrap(err, "create layer store")
	}

	imageStoreBackend, err := image.NewFSStoreBackend(filepath.Join(config.Root, "imagedb"))
	if err != nil {
		return nil, errors.Wrap(err, "create image store backend")
	}

	if m.imageStore, err = image.NewImageStore(imageStoreBackend, m.layerStore); err != nil {
		return nil, errors.Wrap(err, "create image store")
	}

	if m.distributionMetadataStore, err = metadata.NewFSMetadataStore(filepath.Join(config.Root, "distribution")); err != nil {
		return nil, errors.Wrap(err, "create distribution metadata store")
	}

	if m.referenceStore, err = reference.NewReferenceStore(filepath.Join(config.Root, "repositories.json")); err != nil {
		return nil, errors.Wrap(err, "create reference store")
	}

	m.registryService = registry.NewService(config.Registry)
	m.downloadManager = xfer.NewLayerDownloadManager(m.layerStore, config.Transfer.MaxConcurrentDownloads)
	m.uploadManager = xfer.NewLayerUploadManager(config.Transfer.MaxConcurrentUploads)

	return &m, nil
}

func (m *Manager) setupLayerStore(config *Config) error {
	var err error

	options := layer.StoreOptions{
		StorePath:                 config.Root,
		MetadataStorePathTemplate: filepath.Join(config.Root, "image", "%s", "layerdb"),
		GraphDriver:               config.Graph.Driver,
		GraphDriverOptions:        config.Graph.Options,
	}

	if config.RemapRoot.Username != "" && config.RemapRoot.Username != "root" {
		if options.UIDMaps, options.GIDMaps, err = idtools.CreateIDMappings(config.RemapRoot.Username, config.RemapRoot.Group); err != nil {
			return err
		}
	}

	m.layerStore, err = layer.NewStoreFromOptions(options)
	return err
}
