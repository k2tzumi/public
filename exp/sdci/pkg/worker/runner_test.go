package worker

import (
	"context"
	"testing"

	"cirello.io/exp/sdci/pkg/grpc/api"
)

func TestRun(t *testing.T) {
	recipe := &api.Recipe{
		Environment: "RECIPE_MSG=world",
		Commands:    "echo Hello, $RECIPE_MSG;",
	}
	response, err := run(context.Background(), recipe, ".", ".")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", response)
	if response != "Hello, world\n" {
		t.Errorf("unexpected output: %v", response)
	}
}
