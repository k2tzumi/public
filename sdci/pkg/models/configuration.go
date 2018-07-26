package models

import (
	"io"

	"cirello.io/errors"
	yaml "gopkg.in/yaml.v2"
)

// Configuration defines the internal parameters for the application.
type Configuration map[string]Recipe

// LoadConfiguration loads a given fd with YAML content into Configuration.
func LoadConfiguration(r io.Reader) (Configuration, error) {
	var c Configuration
	err := yaml.NewDecoder(r).Decode(&c)
	return c, errors.E(err, "cannot parse configuration")
}
