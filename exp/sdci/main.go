// Command sdci implements a simple and dirty CI service.
package main // import "cirello.io/exp/sdci"
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/user"
	"path/filepath"

	"cirello.io/errors"
	"cirello.io/exp/sdci/git"
	"cirello.io/exp/sdci/runner"
	"github.com/davecgh/go-spew/spew"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("sdci: ")
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalln("cannot load current user information:", err)
	}
	buildsDir := filepath.Join(currentUser.HomeDir, ".sdci", "builds", "src", "github.com")
	if err := os.MkdirAll(buildsDir, os.ModePerm&0700); err != nil {
		log.Fatalln("cannot create .sdci:", err)
	}
	jobs := make(chan *buildRequest, 10)
	go worker(buildsDir, jobs)
	recipes, err := loadRecipes()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/github-webhook/", func(w http.ResponseWriter, r *http.Request) {
		var payload githubHookPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Println("cannot decode payload:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		recipe, ok := recipes[payload.Repository.FullName]
		if !ok {
			log.Println("cannot find recipe for", payload.Repository.FullName)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		jobs <- &buildRequest{
			repoFullName: payload.Repository.FullName,
			commit:       payload.Commit,
			recipe:       recipe,
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		out, err := httputil.DumpRequest(r, true)
		spew.Dump(err)
		fmt.Println(string(out))
	})
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func loadRecipes() (map[string]*runner.Recipe, error) {
	fd, err := os.Open("sdci-config.yaml")
	if err != nil {
		return nil, errors.E(err, "cannot open configuration file")
	}
	var recipes map[string]*runner.Recipe
	err = yaml.NewDecoder(fd).Decode(&recipes)
	return recipes, errors.E(err, "cannot parse configuration")
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

type buildRequest struct {
	repoFullName string
	commit       string
	recipe       *runner.Recipe
}

func worker(buildsDir string, jobs chan *buildRequest) {
	for j := range jobs {
		log.Println("checking out code...")
		dir, name := filepath.Split(j.repoFullName)
		repoDir := filepath.Join(buildsDir, dir, name)
		if err := git.Checkout(j.recipe.Clone, repoDir, j.commit); err != nil {
			log.Println("cannot checkout code:", err)
			continue
		}
		log.Println("building...")
		err := runner.Run(context.Background(), j.recipe, repoDir)
		log.Println("building result:", err)
	}
}
