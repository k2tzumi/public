// Command sdci implements a simple and dirty CI service.
package main // import "cirello.io/exp/sdci"
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"cirello.io/errors"
	"cirello.io/exp/sdci/runner"
	"github.com/davecgh/go-spew/spew"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("sdci: ")
	recipes, err := loadRecipes()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		out, err := httputil.DumpRequest(r, true)
		spew.Dump(err)
		fmt.Println(string(out))
	})
	http.HandleFunc("/github-webhook/", func(w http.ResponseWriter, r *http.Request) {
		var payload githubHookPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		spew.Dump(payload, err)
		recipe := recipes[payload.Repository.FullName]
		spew.Dump(recipe)
		out, err := runner.Run(r.Context(), recipe)
		spew.Dump(out, err)
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
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Sender struct {
		Login     string
		AvatarURL string `json:"avatar_url"`
	}
}
