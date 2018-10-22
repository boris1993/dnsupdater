// Package utils provides utilities about IP addresses, DNS records, and config files.
package utils

import (
	"github.com/boris1993/dnsupdater/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// ReadConfig reads the config.yaml and saves the properties in a variable.
//
// path is the absolute or relative path to the config file.
func ReadConfig(path string) (model.Config, error) {
	var conf model.Config

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(bytes, &conf)

	if err != nil {
		return conf, err
	}

	return conf, nil
}
