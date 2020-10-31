package main

import (
	"errors"
	"flag"
	"github.com/boris1993/dnsupdater/internal/configs"
	"github.com/boris1993/dnsupdater/internal/constants"
	"github.com/boris1993/dnsupdater/internal/helper/aliyun"
	"github.com/boris1993/dnsupdater/internal/helper/cloudflare"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	var err error

	var config = configs.Get()

	// Fetch the current external IP address.
	ipAddress, err := getCurrentIPAddress(config)
	if err != nil {
		log.Fatalln(err)
	}

	// Process CloudFlare records
	err = cloudflare.ProcessRecords(config, ipAddress)
	if err != nil {
		log.Errorln(err)
	}

	err = aliyun.ProcessRecords(config, ipAddress)
	if err != nil {
		log.Errorln(err)
	}

	os.Exit(0)
}

func init() {
	flag.StringVar(&configs.Path, "config", "", "Path to the config file.")
	flag.BoolVar(&configs.Debug, "debug", false, "Enable debug logging.")

	flag.Parse()

	log.SetFormatter(&log.TextFormatter{DisableLevelTruncation: true})

	if configs.Debug == true {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

// getCurrentIPAddress returns the external IP address for your network
func getCurrentIPAddress(config configs.Config) (string, error) {
	if config.System.IPAddrAPI == "" {
		return "", errors.New(constants.ErrIPAddressFetchingAPIEmpty)
	}

	log.Println(constants.MsgCheckingCurrentIPAddr)

	//region fetch your IPv4 address
	resp, err := http.Get(config.System.IPAddrAPI)
	if err != nil {
		return "", err
	}

	// Handle errors when closing the HTTP connection
	defer func() {
		err := resp.Body.Close()

		if err != nil {
			log.Errorln(constants.ErrCloseHTTPConnectionFail, err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Body only contains the IP address
	ipAddress := string(body)
	//endregion

	log.Println(constants.MsgHeaderCurrentIPAddr, ipAddress)

	return ipAddress, nil
}
