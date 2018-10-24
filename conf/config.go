// Package conf provides all models needed by this programme.
package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

var once sync.Once

var Path string
var conf *config

// system describes the system properties in config.yaml
type system struct {
	IPAddrAPI string `yaml:"IPAddrAPI"`
}

// cloudFlare describes the CloudFlare specialised properties in config.yaml
type cloudFlare struct {
	APIEndpoint string `yaml:"APIEndpoint"`
	APIKey      string `yaml:"APIKey"`
	ZoneID      string `yaml:"ZoneID"`
	AuthEmail   string `yaml:"AuthEmail"`
	DomainName  string `yaml:"DomainName"`
}

// config describes the top-level properties in config.yaml
type config struct {
	System     system     `yaml:"System"`
	CloudFlare cloudFlare `yaml:"CloudFlare"`
}

func Get() *config {
	once.Do(func() {
		initConfig()
	})
	return conf
}

// ReadConfig reads the config.yaml and saves the properties in a variable.
//
// path is the absolute or relative path to the config file.
func initConfig() error {
	log.Printf("Loading configuraton from %s.\n", Path)

	bytes, err := ioutil.ReadFile(Path)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, &conf)

	if err != nil {
		return err
	}

	return nil
}
