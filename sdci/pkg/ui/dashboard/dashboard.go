// Package dashboard implements a build web dashboard.
package dashboard // import "cirello.io/exp/sdci/pkg/ui/dashboard"

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/models"
)

// Server implements the dashboard.
type Server struct {
	buildDAO models.BuildRepository
}

// New creates a new builds dashboard.
func New(buildDAO models.BuildRepository) *Server {
	return &Server{
		buildDAO: buildDAO,
	}
}

// ServeContext exposes the build dashboard.
func (s *Server) ServeContext(ctx context.Context, l net.Listener) error {
	srv := &http.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/dashboard/", s.listBuilds)
	go func() {
		select {
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			srv.Shutdown(ctx)
		}
	}()
	srv.Handler = mux
	go func() {
		select {
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			srv.Shutdown(ctx)
		}
	}()
	return errors.E(srv.Serve(l), "error when serving dashboard interface")
}

func (s *Server) listBuilds(w http.ResponseWriter, r *http.Request) {
	repoFullName := strings.TrimPrefix(r.RequestURI, "/dashboard/")
	builds, err := s.buildDAO.ListByRepoFullName(repoFullName)
	if err != nil && errors.RootCause(err) != sql.ErrNoRows {
		log.Println("cannot load builds:", err)
	}
	if err := json.NewEncoder(w).Encode(builds); err != nil {
		log.Println("cannot encode builds:", err)
	}
}
