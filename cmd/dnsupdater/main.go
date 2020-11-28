package main

import (
	"errors"
	"flag"
	"github.com/boris1993/dnsupdater/internal/common"
	"github.com/boris1993/dnsupdater/internal/helper/aliyun"
	"github.com/boris1993/dnsupdater/internal/helper/cloudflare"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	var err error

	config, err := common.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	//region Fetch the current external IP address.
	currentIPv4Address, err := getCurrentIPv4Address(*config)
	if err != nil {
		log.Fatalln(err)
	}

	// For those who doesn't have IPv6 internet access,
	// the currentIPv6Address will be an empty string
	var currentIPv6Address = ""
	if config.System.IPv6AddrAPI == "" {
		log.Info(common.MsgIPv6Disabled)
	} else {
		currentIPv6Address, err = getCurrentIPv6Address(*config)
		if err != nil {
			log.Warnln("Failed to retrieve your IPv6 address.")
			log.Warnln("If you saw the \"no route to host\" error, " +
				"Please verify if you have IPv6 internet access, " +
				"or you can disable IPv6 support by removing the \"IPv6AddrAPI\" in config.yaml.")
			log.Fatalln(err)
		}
	}
	//endregion

	// Process CloudFlare DNS records
	err = cloudflare.ProcessRecords(currentIPv4Address, currentIPv6Address)
	if err != nil {
		log.Errorln(err)
	}

	// Process Aliyun DNS records
	err = aliyun.ProcessRecords(currentIPv4Address, currentIPv6Address)
	if err != nil {
		log.Errorln(err)
	}

	os.Exit(0)
}

func init() {
	flag.StringVar(&common.ConfigFilePath, "config", "", "Path to the config file.")
	flag.BoolVar(&common.Debug, "debug", false, "Enable debug logging.")

	flag.Parse()

	log.SetFormatter(&log.TextFormatter{DisableLevelTruncation: true})

	if common.Debug == true {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

// getCurrentIPv4Address returns the external IP address for your network.
func getCurrentIPv4Address(config common.Config) (string, error) {
	if config.System.IPAddrAPI == "" {
		return "", errors.New(common.ErrIPAddressFetchingAPIEmpty)
	}

	log.Println(common.MsgCheckingCurrentIPv4Addr)

	//region fetch your IPv4 address
	resp, err := http.Get(config.System.IPAddrAPI)
	if err != nil {
		return "", err
	}

	// Handle errors when closing the HTTP connection
	defer func() {
		err := resp.Body.Close()

		if err != nil {
			log.Errorln(common.ErrCloseHTTPConnectionFail, err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Body only contains the IP address
	ipAddress := string(body)
	//endregion

	log.Println(common.MsgHeaderCurrentIPv4Addr, ipAddress)

	return ipAddress, nil
}

// getCurrentIPv6Address returns the external IPv6 address for your network.
// Typically this should be your "temporary" IPv6 address.
func getCurrentIPv6Address(config common.Config) (string, error) {
	if config.System.IPv6AddrAPI == "" {
		return "", errors.New(common.ErrIPAddressFetchingAPIEmpty)
	}

	log.Println(common.MsgCheckingCurrentIPv6Addr)

	resp, err := http.Get(config.System.IPv6AddrAPI)
	if err != nil {
		return "", err
	}

	// Handle errors when closing the HTTP connection
	defer func() {
		err := resp.Body.Close()

		if err != nil {
			log.Errorln(common.ErrCloseHTTPConnectionFail, err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Body only contains the IP address
	ipv6Address := string(body)

	log.Println(common.MsgHeaderCurrentIPv6Addr, ipv6Address)

	return ipv6Address, nil
}
