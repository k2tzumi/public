package web // import "cirello.io/exp/sdci/pkg/web"
import (
	"context"
	"database/sql"
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
	coordinator *coordinator.Coordinator
}

// New creates a new web-facing server.
func New(c *coordinator.Coordinator) *Server {
	return &Server{
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
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		s.coordinator.Enqueue(&models.Build{
			RepoFullName:  payload.Repository.FullName,
			CommitHash:    payload.CommitHash,
			CommitMessage: payload.HeadCommit.Message,
		})
		if err := s.coordinator.Error(); err != nil {
			// TODO: should it really wait forever for shutdown?
			srv.Shutdown(context.Background())
		}
	})
	mux.HandleFunc("/badge/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml;charset=utf-8")
		build, err := s.coordinator.GetLastBuild("ucirello/public")
		if err != nil && errors.RootCause(err) != sql.ErrNoRows {
			log.Println("cannot load last build for repository:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		badge := badgeUnknown
		switch build.Status() {
		case models.Success:
			badge = badgePassing
		case models.Failed:
			badge = badgeFailing
		case models.InProgress:
			badge = badgeRunning
		}
		fmt.Fprint(w, badge)
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

const badgePassing = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="88" height="20">
<linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient>
<clipPath id="a"><rect width="88" height="20" rx="3" fill="#fff"/></clipPath>
<g clip-path="url(#a)"><path fill="#555" d="M0 0h37v20H0z"/><path fill="#4c1" d="M37 0h51v20H37z"/><path fill="url(#b)" d="M0 0h88v20H0z"/></g>
<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">
<text x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="270">build</text>
<text x="195" y="140" transform="scale(.1)" textLength="270">build</text>
<text x="615" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="410">passing</text>
<text x="615" y="140" transform="scale(.1)" textLength="410">passing</text>
</g>
</svg>`

const badgeFailing = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="80" height="20">
<linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient>
<clipPath id="a"><rect width="80" height="20" rx="3" fill="#fff"/></clipPath>
<g clip-path="url(#a)"><path fill="#555" d="M0 0h37v20H0z"/><path fill="#e05d44" d="M37 0h43v20H37z"/><path fill="url(#b)" d="M0 0h80v20H0z"/></g>
<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">
<text x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="270">build</text>
<text x="195" y="140" transform="scale(.1)" textLength="270">build</text>
<text x="575" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="330">failing</text>
<text x="575" y="140" transform="scale(.1)" textLength="330">failing</text>
</g>
</svg>`

const badgeRunning = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="90" height="20">
<linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient>
<clipPath id="a"><rect width="90" height="20" rx="3" fill="#fff"/></clipPath>
<g clip-path="url(#a)"><path fill="#555" d="M0 0h37v20H0z"/><path fill="#9f9f9f" d="M37 0h53v20H37z"/><path fill="url(#b)" d="M0 0h90v20H0z"/></g>
<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">
<text x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="270">build</text>
<text x="195" y="140" transform="scale(.1)" textLength="270">build</text>
<text x="625" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="430">running</text>
<text x="625" y="140" transform="scale(.1)" textLength="430">running</text>
</g>
</svg>`

const badgeUnknown = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="90" height="20">
<linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient>
<clipPath id="a"><rect width="90" height="20" rx="3" fill="#fff"/></clipPath>
<g clip-path="url(#a)"><path fill="#555" d="M0 0h37v20H0z"/><path fill="#9f9f9f" d="M37 0h53v20H37z"/><path fill="url(#b)" d="M0 0h90v20H0z"/></g>
<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">
<text x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="270">build</text>
<text x="195" y="140" transform="scale(.1)" textLength="270">build</text>
<text x="625" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="430">unknown</text>
<text x="625" y="140" transform="scale(.1)" textLength="430">unknown</text>
</g>
</svg>`
