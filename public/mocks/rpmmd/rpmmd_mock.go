package rpmmd_mock

import (
	"github.com/ondrejbudai/osbuild-composer-public/public/store"
	"github.com/ondrejbudai/osbuild-composer-public/public/worker"
)

type Fixture struct {
	StoreFixture *store.Fixture
	Workers      *worker.Server
}
