package runner

import (
	"context"
	"testing"

	"cirello.io/exp/cdci/api"
	"github.com/davecgh/go-spew/spew"
)

func TestRun(t *testing.T) {
	recipe := &api.Recipe{
		Environment: []string{"RECIPE_MSG=hello"},
		Steps: []*api.Step{{
			Environment: []string{"STEP_MSG=world"},
			Commands:    []string{"export"},
		}},
	}
	response := Run(context.TODO(), recipe)
	spew.Dump(response)

}
