package main

import (
	"errors"
	"flag"
	"github.com/boris1993/dnsupdater/alidnsutil"
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
	ipAddress, ipv6Address, err := getCurrentIPAddress(config)

	if err != nil {
		log.Fatalln(err)
	}

	// Process CloudFlare records
	err = cfutil.ProcessRecords(config, ipAddress, ipv6Address)

	if err != nil {
		log.Errorln(err)
	}

	alidnsutil.ProcessRecords(config, ipAddress)

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
func getCurrentIPAddress(config conf.Config) (string, string, error) {
	if config.System.IPAddrAPI == "" || config.System.IPv6AddrAPI == "" {
		return "", "", errors.New(constants.ErrIPAddressFetchingAPIEmpty)
	}

	log.Println(constants.MsgCheckingCurrentIPAddr)

	//region fetch your IPv4 address
	resp, err := http.Get(config.System.IPAddrAPI)

	if err != nil {
		return "", "", err
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
		return "", "", err
	}

	// Body only contains the IP address
	ipAddress := string(body)
	//endregion

	//region fetch your IPv6 address
	resp, err = http.Get(config.System.IPv6AddrAPI)

	if err != nil {
		return "", "", err
	}

	defer func() {
		err := resp.Body.Close()

		if err != nil {
			log.Errorln(constants.ErrCloseHTTPConnectionFail, err)
		}
	}()

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", "", err
	}

	ipv6Address := string(body)
	//endregion

	log.Println(constants.MsgHeaderCurrentIPAddr, ipAddress)
	log.Println(constants.MsgHeaderCurrentIPv6Addr, ipv6Address)

	return ipAddress, ipv6Address, nil
}
