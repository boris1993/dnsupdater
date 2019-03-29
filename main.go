package main

import (
	"errors"
	"flag"
	"github.com/boris1993/dnsupdater/cfutil"
	"github.com/boris1993/dnsupdater/conf"
	"github.com/boris1993/dnsupdater/constants"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	var err error

	var config = conf.Get()

	// Fetch the current external IP address.
	currentIPAddress, err := getCurrentIPAddress(config)

	if err != nil {
		log.Fatalln(err)
	}

	// Process CloudFlare records
	err = cfutil.ProcessRecords(config, currentIPAddress)

	if err != nil {
		log.Errorln(err)
	}

	os.Exit(0)
}

func init() {
	flag.StringVar(&conf.Path, "config", "", "Path to the config file.")
	flag.BoolVar(&conf.Debug, "debug", false, "Enable debug logging.")

	flag.Parse()

	log.SetFormatter(&log.TextFormatter{DisableLevelTruncation: true})

	if conf.Debug == true {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

// getCurrentIPAddress returns the external IP address for your network
func getCurrentIPAddress(config *conf.Config) (string, error) {
	if config.System.IPAddrAPI == "" {
		return "", errors.New(constants.ErrIPAddressFetchingAPIEmpty)
	}

	log.Println(constants.MsgCheckingCurrentIPAddr)

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

	log.Println(constants.MsgHeaderCurrentIPAddr, ipAddress)

	return ipAddress, nil
}
