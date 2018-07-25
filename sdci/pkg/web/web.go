package web // import "cirello.io/exp/sdci/pkg/web"
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/coordinator"
	"cirello.io/exp/sdci/pkg/models"
)

// Server implements the web-facing part of the CI service. For now, compatible
// only with Github Webhooks.
type Server struct {
	recipes     map[string]*models.Recipe // map of fullName to recipes
	coordinator *coordinator.Coordinator
}

// New creates a new web-facing server.
func New(r map[string]*models.Recipe, c *coordinator.Coordinator) *Server {
	return &Server{
		recipes:     r,
		coordinator: c,
	}
}

// Serve handles the HTTP requests.
func (s *Server) Serve(l net.Listener) error {
	srv := &http.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/github-webhook/", func(w http.ResponseWriter, r *http.Request) {
		var payload githubHookPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Println("cannot decode payload:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		recipe, ok := s.recipes[payload.Repository.FullName]
		if !ok {
			log.Println("cannot find recipe for", payload.Repository.FullName)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		s.coordinator.Enqueue(&models.Build{
			RepoFullName:  payload.Repository.FullName,
			CommitHash:    payload.CommitHash,
			CommitMessage: payload.HeadCommit.Message,
			Recipe:        recipe,
		})
		if err := s.coordinator.Error(); err != nil {
			// TODO: should it really wait forever for shutdown?
			srv.Shutdown(context.Background())
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		out, err := httputil.DumpRequest(r, true)
		fmt.Println(string(out), err)
	})
	srv.Handler = mux
	return errors.E(srv.Serve(l), "error when serving web interface")
}

type githubHookPayload struct {
	CommitHash string `json:"after"`
	HeadCommit struct {
		Message string `json:"message"`
	} `json:"head_commit"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Sender struct {
		Login     string
		AvatarURL string `json:"avatar_url"`
	}
}
