package models

import (
	"io"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/coordinator/api"
	yaml "gopkg.in/yaml.v2"
)

// Configuration defines the internal parameters for the application.
type Configuration map[string]api.Recipe

// LoadConfiguration loads a given fd with YAML content into Configuration.
func LoadConfiguration(r io.Reader) (Configuration, error) {
	var c Configuration
	err := yaml.NewDecoder(r).Decode(&c)
	return c, errors.E(err, "cannot parse configuration")
}
