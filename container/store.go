package container

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/engine-api/types"
	"github.com/docker/swarmkit/ioutils"
	"github.com/pkg/errors"
)

type Store struct {
	root   string
	byId   string
	byName string
}

type Config struct {
	Root string
}

func NewStore(config *Config) (*Store, error) {
	s := &Store{
		root:   config.Root,
		byId:   filepath.Join(config.Root, "id"),
		byName: filepath.Join(config.Root, "name"),
	}

	for _, p := range []string{s.byId, s.byName} {
		if err := os.MkdirAll(p, 0700); err != nil && !os.IsExist(err) {
			return nil, err
		}
	}

	return s, nil
}

func (s *Store) Register(container types.ContainerCreateConfig, imageId string) (string, error) {
	var err error
	id := stringid.GenerateNonCryptoID()
	name := container.Name

	defer func() {
		if err != nil {
			s.deleteContainer(id, name)
		}
	}()

	if name == "" {
		name, err = s.generateName(id)
		if err != nil {
			return "", err
		}
	} else {
		// err is local so that defer won't see it
		if err := s.recordName(name, id); err != nil {
			return "", errors.Wrap(err, "name in use")
		}
	}

	if err = os.MkdirAll(s.idPath(id), 0700); err != nil {
		return "", err
	}

	if err := s.saveFile(id, "name", []byte(name)); err != nil {
		return "", err
	}

	if err = s.saveJsonFile(id, "hostconfig.json", container.HostConfig); err != nil {
		return "", err
	}

	if err = s.saveJsonFile(id, "config.json", container.Config); err != nil {
		return "", err
	}

	if err = s.saveJsonFile(id, "networkconfig.json", container.NetworkingConfig); err != nil {
		return "", err
	}

	if err = os.MkdirAll(path.Join(s.idPath(id), "image"), 0700); err != nil && !os.IsExist(err) {
		return "", err
	}

	f, err := os.Create(path.Join(s.idPath(id), "image", imageId))
	if err != nil {
		return "", err
	}
	f.Close()

	return id, nil
}

func (s *Store) saveFile(id, file string, content []byte) error {
	return ioutils.AtomicWriteFile(path.Join(s.idPath(id), file), content, 0700)
}

func (s *Store) saveJsonFile(id, file string, obj interface{}) error {
	if obj == nil {
		return nil
	}
	content, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return s.saveFile(id, file, content)
}

func (s *Store) deleteContainer(id, name string) error {
	s.releaseName(name)
	if id == "" {
		return nil
	}
	return os.RemoveAll(s.idPath(id))
}

func (s *Store) generateName(id string) (string, error) {
	for i := 0; i < 6; i++ {
		name := namesgenerator.GetRandomName(i)
		if s.recordName(name, id) == nil {
			return name, nil
		}
	}
	return id[:12], s.recordName(id[:12], id)
}

func (s *Store) idPath(id string) string {
	return filepath.Join(s.byId, id)
}

func (s *Store) namePath(name string) string {
	return filepath.Join(s.byName, name)
}

func (s *Store) releaseName(name string) error {
	if name == "" {
		return nil
	}
	err := os.Remove(s.namePath(name))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *Store) recordName(name, id string) error {
	if err := os.Symlink(s.idPath(id), s.namePath(name)); err != nil {
		return errors.Wrap(err, "record name")
	}
	return nil
}
