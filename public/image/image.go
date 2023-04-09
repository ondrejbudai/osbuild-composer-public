package image

import (
	"math/rand"

	"github.com/ondrejbudai/osbuild-composer-public/public/artifact"
	"github.com/ondrejbudai/osbuild-composer-public/public/manifest"
	"github.com/ondrejbudai/osbuild-composer-public/public/rpmmd"
	"github.com/ondrejbudai/osbuild-composer-public/public/runner"
)

type ImageKind interface {
	Name() string
	InstantiateManifest(m *manifest.Manifest, repos []rpmmd.RepoConfig, runner runner.Runner, rng *rand.Rand) (*artifact.Artifact, error)
}

type Base struct {
	name string
}

func (img Base) Name() string {
	return img.name
}

func NewBase(name string) Base {
	return Base{
		name: name,
	}
}
