// Package conf provides all models needed by this programme.
package conf

import (
	"github.com/boris1993/dnsupdater/constants"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var once sync.Once

var Debug bool

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
		err := initConfig()

		if err != nil {
			log.Fatalln(err)
		}
	})
	return conf
}

// initConfig reads the config.yaml and saves the properties in a variable.
func initConfig() error {
	if Path == "" {
		absPath, err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			log.Fatalln(err)
		}

		Path = filepath.Join(absPath, "config.yaml")
	}

	log.Println(constants.MsgHeaderLoadingConfig, Path)

	bytes, err := ioutil.ReadFile(Path)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, &conf)

	if err != nil {
		return err
	}

	if Debug {
		printDebugInfo()
	}

	return nil
}

// printDebugInfo prints the configurations loaded from the file.
func printDebugInfo() {
	log.Debugf("%15v: %s", "IPAddrAPI", conf.System.IPAddrAPI)
	log.Debugf("%15v: %s", "APIEndpoint", conf.CloudFlare.APIEndpoint)
	log.Debugf("%15v: %s", "APIKey", conf.CloudFlare.APIKey)
	log.Debugf("%15v: %s", "ZoneID", conf.CloudFlare.ZoneID)
	log.Debugf("%15v: %s", "AuthEmail", conf.CloudFlare.AuthEmail)
	log.Debugf("%15v: %s", "DomainName", conf.CloudFlare.DomainName)
}
