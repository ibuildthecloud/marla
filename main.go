package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/reexec"
	"github.com/ibuildthecloud/marla/container"
	"github.com/ibuildthecloud/marla/daemon"
	"github.com/ibuildthecloud/marla/rootfs/docker"
	"github.com/ibuildthecloud/marla/server"

	// Register overlay
	_ "github.com/docker/docker/daemon/graphdriver/overlay"
	// Setup docker-untar, docker-applyLayer
	_ "github.com/docker/docker/pkg/chrootarchive"
)

func main() {
	if reexec.Init() {
		return
	}

	if err := mainWithError(); err != nil {
		logrus.Fatal(err)
	}
}

func mainWithError() error {
	logrus.SetLevel(logrus.DebugLevel)

	if err := SetUmask(); err != nil {
		logrus.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cwd = filepath.Join(cwd, "runtime")
	root := filepath.Join(cwd, "overlay")

	// work around docker bug
	os.MkdirAll(filepath.Join(root, "image"), 0700)

	daemon, err := daemon.New(&daemon.Config{
		Root: cwd,
		DockerRootFS: docker.Config{
			Root: filepath.Join(root, "image"),
			Graph: docker.GraphConfig{
				Driver: "overlay",
			},
			Transfer: docker.TransferConfig{
				MaxConcurrentUploads:   5,
				MaxConcurrentDownloads: 5,
			},
		},
		Container: container.Config{
			Root: filepath.Join(root, "container"),
		},
	})
	if err != nil {
		logrus.Fatal(err)
	}

	server, err := server.New(daemon)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := server.ListenAndServe(); err != nil {
		logrus.Fatal(err)
	}

	return nil
}

// TODO: move to pkg/umask
func SetUmask() error {
	desiredUmask := 0022
	syscall.Umask(desiredUmask)
	if umask := syscall.Umask(desiredUmask); umask != desiredUmask {
		return fmt.Errorf("failed to set umask: expected %#o, got %#o", desiredUmask, umask)
	}

	return nil
}
