package models // import "cirello.io/exp/sdci/pkg/models"

// Recipe defines the execution steps and environment.
type Recipe struct {
	Concurrency  int    `db:"-" yaml:"concurrency" json:"concurrency,omitempty"`
	Clone        string `db:"-" yaml:"clone" json:"clone,omitempty"`
	SlackWebhook string `db:"-" yaml:"slack_webhook" json:"slack_webhook,omitempty"`
	GithubSecret string `db:"-" yaml:"github_secret" json:"github_secret,omitempty"`
	Environment  string `db:"environment" yaml:"environment" json:"environment,omitempty"`
	Commands     string `db:"commands" yaml:"commands" json:"commands,omitempty"`
}
