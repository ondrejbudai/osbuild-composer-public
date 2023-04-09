package distro_mock

import (
	"github.com/ondrejbudai/osbuild-composer-public/public/distro/test_distro"
	"github.com/ondrejbudai/osbuild-composer-public/public/distroregistry"
)

func NewDefaultRegistry() (*distroregistry.Registry, error) {
	testDistro := test_distro.New()
	if testDistro == nil {
		panic("Attempt to register test distro failed")
	}
	return distroregistry.New(nil, testDistro)
}
