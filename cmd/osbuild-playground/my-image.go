package main

import (
	"math/rand"

	"github.com/ondrejbudai/osbuild-composer-public/public/artifact"
	"github.com/ondrejbudai/osbuild-composer-public/public/disk"
	"github.com/ondrejbudai/osbuild-composer-public/public/manifest"
	"github.com/ondrejbudai/osbuild-composer-public/public/platform"
	"github.com/ondrejbudai/osbuild-composer-public/public/rpmmd"
	"github.com/ondrejbudai/osbuild-composer-public/public/runner"
)

type MyImage struct {
	MyOption string `json:"my_option"`
}

func (img *MyImage) Name() string {
	return "my-image"
}

func init() {
	AddImageType(&MyImage{})
}

func (img *MyImage) InstantiateManifest(m *manifest.Manifest,
	repos []rpmmd.RepoConfig,
	runner runner.Runner,
	rng *rand.Rand) (*artifact.Artifact, error) {
	// Let's create a simple raw image!

	// configure a build pipeline
	build := manifest.NewBuild(m, runner, repos)
	build.Checkpoint()

	// create an x86_64 platform with bios boot
	platform := &platform.X86{
		BIOS: true,
	}

	// TODO: add helper
	pt, err := disk.NewPartitionTable(&basePT, nil, 0, false, nil, rng)
	if err != nil {
		panic(err)
	}

	// create a minimal bootable OS tree
	os := manifest.NewOS(m, build, platform, repos)
	os.PartitionTable = pt   // we need a partition table
	os.KernelName = "kernel" // use the default fedora kernel

	// create a raw image containing the OS tree created above
	raw := manifest.NewRawImage(m, build, os)
	artifact := raw.Export()

	return artifact, nil
}
