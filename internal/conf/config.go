// Package conf provides all models needed by this programme.
package conf

import (
	"encoding/json"
	"github.com/boris1993/dnsupdater/internal/common"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"sync"
)

var once = new(sync.Once)

var Debug bool

var ConfigFilePath string
var conf Config
var errorInInitConfig error

type ExternalAddressResponseType string

const (
	Text ExternalAddressResponseType = "text"
	Json ExternalAddressResponseType = "json"
)

// Config describes the top-level properties in Config.yaml
type Config struct {
	System            System       `yaml:"System"`
	CloudFlareRecords []CloudFlare `yaml:"CloudFlareRecords"`
	AliDNSRecords     []AliDNS     `yaml:"AliDNSRecords"`
}

// System describes the System properties in Config.yaml
type System struct {
	IPv4      PublicIPAddressEndpointConfig `yaml:"IPv4"`
	IPv6      PublicIPAddressEndpointConfig `yaml:"IPv6"`
	Endpoints DNSProviderEndpointConfig     `yaml:"Endpoints"`
}

type PublicIPAddressEndpointConfig struct {
	Enabled        bool                        `yaml:"Enabled"`
	IPAddrAPI      string                      `yaml:"IPAddrAPI"`
	ResponseType   ExternalAddressResponseType `yaml:"ResponseType"`
	IPAddrJsonPath string                      `yaml:"IPAddrJsonPath"`
}

type DNSProviderEndpointConfig struct {
	CloudFlareAPIEndpoint string `yaml:"CloudFlareAPIEndpoint"`
	AliyunAPIEndpoint     string `yaml:"AliyunAPIEndpoint"`
}

// CloudFlare describes the CloudFlare specialised properties in Config.yaml
type CloudFlare struct {
	//APIEndpoint string `yaml:"APIEndpoint"`
	APIKey     string `yaml:"APIKey"`
	ZoneID     string `yaml:"ZoneID"`
	AuthEmail  string `yaml:"AuthEmail"`
	DomainName string `yaml:"DomainName"`
	DomainType string `yaml:"DomainType"`
}

type AliDNS struct {
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	RegionID        string `yaml:"RegionID"`
	DomainName      string `yaml:"DomainName"`
	DomainType      string `yaml:"DomainType"`
}

func GetConfig() (*Config, error) {
	once.Do(func() {
		err := initConfig()

		if err != nil {
			log.Errorln(err)
			errorInInitConfig = err
		}
	})

	if errorInInitConfig != nil {
		return &conf, errorInInitConfig
	}

	return &conf, nil
}

// initConfig reads the Config.yaml and saves the properties in a variable.
func initConfig() error {
	if ConfigFilePath == "" {
		absPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}

		ConfigFilePath = filepath.Join(absPath, "config.yaml")
	}

	log.Println(common.MsgHeaderLoadingConfig, ConfigFilePath)

	bytes, err := os.ReadFile(ConfigFilePath)

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
	bytes, _ := json.Marshal(conf)
	log.Debug(string(bytes))
}
