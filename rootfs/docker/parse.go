package docker

import (
	"github.com/docker/distribution/digest"
	"github.com/docker/docker/reference"
)

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
