package models

import (
	"reflect"
	"strings"
	"testing"

	"cirello.io/exp/sdci/pkg/grpc/api"
)

func TestConfigurationParser(t *testing.T) {
	yaml := strings.NewReader(yamlSample)
	got, err := LoadConfiguration(yaml)
	if err != nil {
		t.Fatalf("cannot parse configuration: %v", err)
	}
	expected := Configuration{
		"org/account": api.Recipe{
			Concurrency:  2,
			Clone:        "git@github.com:org/account.git",
			SlackWebhook: "https://hooks.slack.com/services/AAAA/BBB/CCC",
			GithubSecret: "ghsecret",
			Environment:  "ENV1=1\nENV2=2\n",
			Commands:     "vgo test ./errors/... ./supervisor/...\necho OK\n",
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("wrong parsing. got: %#v\nexpected:%#v\n", got, expected)
	}
}

const yamlSample = `---
org/account:
  concurrency: 2
  clone: git@github.com:org/account.git
  slack_webhook: https://hooks.slack.com/services/AAAA/BBB/CCC
  github_secret: ghsecret
  environment: |
    ENV1=1
    ENV2=2
  commands: |
    vgo test ./errors/... ./supervisor/...
    echo OK
`
