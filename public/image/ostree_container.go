package image

import (
	"math/rand"

	"github.com/ondrejbudai/osbuild-composer-public/public/artifact"
	"github.com/ondrejbudai/osbuild-composer-public/public/environment"
	"github.com/ondrejbudai/osbuild-composer-public/public/manifest"
	"github.com/ondrejbudai/osbuild-composer-public/public/ostree"
	"github.com/ondrejbudai/osbuild-composer-public/public/platform"
	"github.com/ondrejbudai/osbuild-composer-public/public/rpmmd"
	"github.com/ondrejbudai/osbuild-composer-public/public/runner"
	"github.com/ondrejbudai/osbuild-composer-public/public/workload"
)

type OSTreeContainer struct {
	Base
	Platform               platform.Platform
	OSCustomizations       manifest.OSCustomizations
	Environment            environment.Environment
	Workload               workload.Workload
	OSTreeRef              string
	OSTreeParent           *ostree.CommitSpec
	OSVersion              string
	ExtraContainerPackages rpmmd.PackageSet // FIXME: this is never read
	ContainerLanguage      string
	Filename               string
}

func NewOSTreeContainer(ref string) *OSTreeContainer {
	return &OSTreeContainer{
		Base:      NewBase("ostree-container"),
		OSTreeRef: ref,
	}
}

func (img *OSTreeContainer) InstantiateManifest(m *manifest.Manifest,
	repos []rpmmd.RepoConfig,
	runner runner.Runner,
	rng *rand.Rand) (*artifact.Artifact, error) {
	buildPipeline := manifest.NewBuild(m, runner, repos)
	buildPipeline.Checkpoint()

	osPipeline := manifest.NewOS(m, buildPipeline, img.Platform, repos)
	osPipeline.OSCustomizations = img.OSCustomizations
	osPipeline.Environment = img.Environment
	osPipeline.Workload = img.Workload
	osPipeline.OSTreeRef = img.OSTreeRef
	osPipeline.OSTreeParent = img.OSTreeParent

	commitPipeline := manifest.NewOSTreeCommit(m, buildPipeline, osPipeline, img.OSTreeRef)
	commitPipeline.OSVersion = img.OSVersion

	nginxConfigPath := "/etc/nginx.conf"
	listenPort := "8080"

	serverPipeline := manifest.NewOSTreeCommitServer(m,
		buildPipeline,
		img.Platform,
		repos,
		commitPipeline,
		nginxConfigPath,
		listenPort)
	serverPipeline.Language = img.ContainerLanguage

	containerPipeline := manifest.NewOCIContainer(m, buildPipeline, serverPipeline)
	containerPipeline.Cmd = []string{"nginx", "-c", nginxConfigPath}
	containerPipeline.ExposedPorts = []string{listenPort}
	containerPipeline.Filename = img.Filename
	artifact := containerPipeline.Export()

	return artifact, nil
}
