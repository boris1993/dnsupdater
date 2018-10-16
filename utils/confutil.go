package utils

import (
	"github.com/boris1993/dnsupdater/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

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
