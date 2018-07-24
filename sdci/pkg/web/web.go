package web // import "cirello.io/exp/sdci/pkg/web"
import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/coordinator"
)

// Server implements the web-facing part of the CI service. For now, compatible
// only with Github Webhooks.
type Server struct {
	recipes     map[string]*coordinator.Recipe // map of fullName to recipes
	coordinator *coordinator.Coordinator
}

// New creates a new web-facing server.
func New(recipes map[string]*coordinator.Recipe,
	coordinator *coordinator.Coordinator) *Server {
	return &Server{
		recipes:     recipes,
		coordinator: coordinator,
	}
}

// Serve handles the HTTP requests.
func (s *Server) Serve(l net.Listener) error {
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
		s.coordinator.Enqueue(&coordinator.Build{
			RepoFullName: payload.Repository.FullName,
			Commit:       payload.Commit,
			Recipe:       recipe,
		})
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		out, err := httputil.DumpRequest(r, true)
		fmt.Println(string(out), err)
	})
	return errors.E(http.Serve(l, mux), "error when serving web interface")
}

type githubHookPayload struct {
	Commit     string `json:"after"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Sender struct {
		Login     string
		AvatarURL string `json:"avatar_url"`
	}
}
