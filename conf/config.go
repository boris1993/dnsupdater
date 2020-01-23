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
var conf Config

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
}

func Get() Config {
	once.Do(func() {
		err := initConfig()

		if err != nil {
			log.Fatalln(err)
		}
	})
	return conf
}

// initConfig reads the Config.yaml and saves the properties in a variable.
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
	log.Debugf("%21v: %s", "IPAddrAPI", conf.System.IPAddrAPI)
	log.Debugf("%21v: %s", "CloudFlareAPIEndpoint", conf.System.CloudFlareAPIEndpoint)
	log.Debugln()

	for _, item := range conf.CloudFlareRecords {
		log.Debugln("========== CloudFlare DNS Record ==========")
		log.Debugf("%10v: %s", "APIKey", item.APIKey)
		log.Debugf("%10v: %s", "ZoneID", item.ZoneID)
		log.Debugf("%10v: %s", "AuthEmail", item.AuthEmail)
		log.Debugf("%10v: %s", "DomainName", item.DomainName)
		log.Debugf("%10v: %s", "DomainType", item.DomainType)
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
