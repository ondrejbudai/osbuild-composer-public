package rpmmd_mock

import (
	dnfjson_mock "github.com/ondrejbudai/osbuild-composer-public/public/mocks/dnfjson"
	"github.com/ondrejbudai/osbuild-composer-public/public/store"
	"github.com/ondrejbudai/osbuild-composer-public/public/worker"
)

type Fixture struct {
	StoreFixture *store.Fixture
	Workers      *worker.Server
	dnfjson_mock.ResponseGenerator
}
