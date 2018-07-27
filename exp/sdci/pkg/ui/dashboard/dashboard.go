package dashboard // import "cirello.io/exp/sdci/pkg/ui/dashboard"

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/models"
)

// Server implements the dashboard.
type Server struct {
	buildDAO *models.BuildDAO
}

// New creates a new builds dashboard.
func New(buildDAO *models.BuildDAO) *Server {
	return &Server{
		buildDAO: buildDAO,
	}
}

// Serve exposes the build dashboard.
func (s *Server) Serve(l net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/dashboard/", s.listBuilds)
	return http.Serve(l, mux)
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
