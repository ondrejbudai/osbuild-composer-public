package cloudapi

import (
	"net/http"

	"github.com/osbuild/images/pkg/distrofactory"
	"github.com/osbuild/images/pkg/reporegistry"

	v2 "github.com/ondrejbudai/osbuild-composer-public/public/cloudapi/v2"
	"github.com/ondrejbudai/osbuild-composer-public/public/worker"
)

type Server struct {
	v2 *v2.Server
}

func NewServer(workers *worker.Server, distros *distrofactory.Factory, repos *reporegistry.RepoRegistry, config v2.ServerConfig) *Server {
	server := &Server{
		v2: v2.NewServer(workers, distros, repos, config),
	}
	return server
}

func (server *Server) V2(path string) http.Handler {
	return server.v2.Handler(path)
}

func (server *Server) Shutdown() {
	server.v2.Shutdown()
}
