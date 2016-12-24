package docker

import (
	"errors"

	"github.com/docker/distribution/digest"
	"github.com/docker/docker/image"
	"github.com/docker/docker/image/v1"
	"github.com/docker/docker/reference"
)

var errImageNotFound = errors.New("Image not found")

func imageTagToReference(image, tag string) (reference.Named, error) {
	ref, err := reference.ParseNamed(image)
	if err != nil {
		return nil, err
	}

	if tag == "" {
		return ref, nil
	}

	if d, err := digest.ParseDigest(tag); err == nil {
		return reference.WithDigest(ref, d)
	} else {
		return reference.WithTag(ref, tag)
	}
}

func parseId(id string) digest.Digest {
	if err := v1.ValidateID(id); err != nil {
		return ""
	}

	id = "sha256:" + id
	dgst, err := digest.ParseDigest(id)
	if err == nil {
		return dgst
	} else {
		return ""
	}
}

func (m *Manager) resolveImageIDByID(name digest.Digest) (image.ID, error) {
	imageId := image.ID(name)
	if _, err := m.imageStore.Get(imageId); err != nil {
		return "", errImageNotFound
	}
	return imageId, nil
}

func (m *Manager) resolveImageIDByRef(name string) (image.ID, error) {
	ref, err := reference.ParseNamed(name)
	if err != nil {
		return "", errImageNotFound
	}

	// By refernce
	if imageID, err := m.referenceStore.Get(ref); err == nil {
		return imageID, nil
	}

	// TODO: Docker does this but I don't get why, commenting out until I figure out why
	//if tagged, ok := ref.(reference.NamedTagged); ok {
	//	if id, err := m.imageStore.Search(tagged.Tag()); err == nil {
	//		for _, namedRef := range m.referenceStore.References(id) {
	//			if namedRef.Name() == ref.Name() {
	//				return id, nil
	//			}
	//		}
	//	}
	//}

	// Search based on ID
	if id, err := m.imageStore.Search(name); err == nil {
		return id, nil
	}

	return "", errImageNotFound
}

func (m *Manager) ResolveImageID(name string) (image.ID, error) {
	id := parseId(name)
	if id == "" {
		return m.resolveImageIDByRef(name)
	} else {
		return m.resolveImageIDByID(id)
	}
}
