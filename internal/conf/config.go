// Package conf provides all models needed by this programme.
package conf

import (
	"github.com/boris1993/dnsupdater/internal/common"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var once = new(sync.Once)

var Debug bool

var ConfigFilePath string
var conf Config
var errorInInitConfig error

// Config describes the top-level properties in Config.yaml
type Config struct {
	System            System       `yaml:"System"`
	CloudFlareRecords []CloudFlare `yaml:"CloudFlareRecords"`
	AliDNSRecords     []AliDNS     `yaml:"AliDNSRecords"`
}

// System describes the System properties in Config.yaml
type System struct {
	IPAddrAPI             string `yaml:"IPAddrAPI"`
	IPv6AddrAPI           string `yaml:"IPv6AddrAPI"`
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

	bytes, err := ioutil.ReadFile(ConfigFilePath)

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
	log.Debugf("%21v: %s", "IPAddrAPI", conf.System.IPAddrAPI)
	log.Debugf("%21v: %s", "CloudFlareAPIEndpoint", conf.System.CloudFlareAPIEndpoint)
	log.Debugln()

	for _, item := range conf.CloudFlareRecords {
		log.Debugln("========== CloudFlare DNS Record ==========")
		log.Debugf("%10v: %s", "APIKey", item.APIKey)
		log.Debugf("%10v: %s", "ZoneID", item.ZoneID)
		log.Debugf("%10v: %s", "AuthEmail", item.AuthEmail)
		log.Debugf("%10v: %s", "DomainName", item.DomainName)
		log.Debugln()
	}

	for _, item := range conf.AliDNSRecords {
		log.Debugln("========== Aliyun DNS Record ==========")
		log.Debugf("%15v: %s", "AccessID", item.AccessKeyID)
		log.Debugf("%15v: %s", "AccessKeySecret", item.AccessKeySecret)
		log.Debugf("%15v: %s", "RegionID", item.RegionID)
		log.Debugf("%15v: %s", "DomainName", item.DomainName)
		log.Debugln()
	}
}
