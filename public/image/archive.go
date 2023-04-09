package image

import (
	"math/rand"

	"github.com/ondrejbudai/osbuild-composer-public/public/artifact"
	"github.com/ondrejbudai/osbuild-composer-public/public/environment"
	"github.com/ondrejbudai/osbuild-composer-public/public/manifest"
	"github.com/ondrejbudai/osbuild-composer-public/public/platform"
	"github.com/ondrejbudai/osbuild-composer-public/public/rpmmd"
	"github.com/ondrejbudai/osbuild-composer-public/public/runner"
	"github.com/ondrejbudai/osbuild-composer-public/public/workload"
)

type Archive struct {
	Base
	Platform         platform.Platform
	OSCustomizations manifest.OSCustomizations
	Environment      environment.Environment
	Workload         workload.Workload
	Filename         string
}

func NewArchive() *Archive {
	return &Archive{
		Base: NewBase("archive"),
	}
}

func (img *Archive) InstantiateManifest(m *manifest.Manifest,
	repos []rpmmd.RepoConfig,
	runner runner.Runner,
	rng *rand.Rand) (*artifact.Artifact, error) {
	buildPipeline := manifest.NewBuild(m, runner, repos)
	buildPipeline.Checkpoint()

	osPipeline := manifest.NewOS(m, buildPipeline, img.Platform, repos)
	osPipeline.OSCustomizations = img.OSCustomizations
	osPipeline.Environment = img.Environment
	osPipeline.Workload = img.Workload

	tarPipeline := manifest.NewTar(m, buildPipeline, &osPipeline.Base, "archive")
	tarPipeline.Filename = img.Filename
	artifact := tarPipeline.Export()

	return artifact, nil
}
