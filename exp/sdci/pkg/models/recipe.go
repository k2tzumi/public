package models // import "cirello.io/exp/sdci/pkg/models"

// Recipe defines the execution steps and environment.
type Recipe struct {
	Concurrency  int    `db:"-" yaml:"concurrency"`
	Clone        string `db:"-" yaml:"clone"`
	SlackWebhook string `db:"-" yaml:"slack_webhook"`
	Environment  string `db:"environment" yaml:"environment"`
	Commands     string `db:"commands" yaml:"commands"`
}
