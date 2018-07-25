package models // import "cirello.io/exp/sdci/pkg/models"

// Recipe defines the execution steps and environment.
type Recipe struct {
	Clone        string `db:"clone" yaml:"clone"`
	SlackWebhook string `db:"slack_webhook" yaml:"slack_webhook"`
	Environment  string `db:"environment" yaml:"environment"`
	Commands     string `db:"commands" yaml:"commands"`
}
