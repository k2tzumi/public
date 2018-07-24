package runner

import (
	"context"
	"testing"
)

func TestRun(t *testing.T) {
	recipe := &Recipe{
		Environment: []string{"RECIPE_MSG=world"},
		Commands:    "echo Hello, $RECIPE_MSG;",
	}
	response, err := Run(context.TODO(), recipe)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", response)
	if response != "Hello, world\n" {
		t.Errorf("unexpected output: %v", response)
	}
}
