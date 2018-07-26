package models

import (
	"reflect"
	"strings"
	"testing"
)

func TestConfigurationParser(t *testing.T) {
	yaml := strings.NewReader(yamlSample)
	got, err := LoadConfiguration(yaml)
	if err != nil {
		t.Fatalf("cannot parse configuration: %v", err)
	}
	expected := Configuration{
		"org/account": Recipe{
			Concurrency:  2,
			Clone:        "git@github.com:org/account.git",
			SlackWebhook: "https://hooks.slack.com/services/AAAA/BBB/CCC",
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
  environment: |
    ENV1=1
    ENV2=2
  commands: |
    vgo test ./errors/... ./supervisor/...
    echo OK
`