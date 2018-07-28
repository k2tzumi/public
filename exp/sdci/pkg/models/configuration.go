package models

import (
	"io"

	"cirello.io/errors"
	"cirello.io/exp/sdci/pkg/grpc/api"
	yaml "gopkg.in/yaml.v2"
)

// Configuration defines the internal parameters for the application.
type Configuration map[string]api.Recipe

// LoadConfiguration loads a given fd with YAML content into Configuration.
func LoadConfiguration(r io.Reader) (Configuration, error) {
	var c Configuration
	err := yaml.NewDecoder(r).Decode(&c)
	for k, v := range c {
		if v.Concurrency == 0 {
			v.Concurrency = 1
			c[k] = v
		}
	}
	return c, errors.E(err, "cannot parse configuration")
}
