package main

import (
	"errors"
	"flag"
	"github.com/boris1993/dnsupdater/internal/common"
	"github.com/boris1993/dnsupdater/internal/conf"
	"github.com/boris1993/dnsupdater/internal/globals"
	"github.com/boris1993/dnsupdater/internal/helper"
	"github.com/boris1993/dnsupdater/internal/helper/aliyun"
	"github.com/boris1993/dnsupdater/internal/helper/cloudflare"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
)

var config *conf.Config

func main() {
	interfaceImplementationCheck()

	var err error

	config, err = conf.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	//region Fetch the current external IP address.
	globals.IPv4Address, err = getCurrentIPv4Address(*config)
	if err != nil {
		log.Errorln("Error occurred when checking and processing IPv4 records", err)
	}

	// For those who doesn't have IPv6 internet access,
	// the currentIPv6Address will be an empty string
	globals.IPv6Address, err = getCurrentIPv6Address(*config)
	if err != nil {
		log.Warnln("Failed to retrieve your IPv6 address.")
		log.Warnln("If you saw the \"no route to host\" error, " +
			"Please verify if you have IPv6 internet access, " +
			"or you can disable IPv6 support by removing the \"IPv6AddrAPI\" in config.yaml.")
		log.Errorln(err)
	}
	//endregion

	handlers := []helper.DDNSHelperInterface{
		cloudflare.CloudFlareDDNSHandler{},
		aliyun.AliyunDDNSHandler{},
	}

	for _, handler := range handlers {
		err = handler.ProcessRecords(globals.IPv4Address, globals.IPv6Address)
		if err != nil {
			log.Errorln(err)
		}
	}

	os.Exit(0)
}

func interfaceImplementationCheck() {
	var _ helper.DDNSHelperInterface = &cloudflare.CloudFlareDDNSHandler{}
	var _ helper.DDNSHelperInterface = &aliyun.AliyunDDNSHandler{}
}

// getCurrentIPv4Address returns the external IP address for your network.
func getCurrentIPv4Address(config conf.Config) (string, error) {
	if config.System.IPv4.Enabled == false {
		log.Infoln(common.MsgIPv4Disabled)
		return "", nil
	}

	log.Println(common.MsgCheckingCurrentIPv4Addr)

	if config.System.IPv4.IPAddrAPI == "" {
		return "", errors.New(common.ErrIPAddressFetchingAPIEmpty)
	}

	//region fetch your IPv4 address
	resp, err := http.Get(config.System.IPv4.IPAddrAPI)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Body only contains the IP address
	ipAddress := strings.TrimSpace(string(body))
	//endregion

	log.Println(common.MsgHeaderCurrentIPv4Addr, ipAddress)

	return ipAddress, nil
}

// getCurrentIPv6Address returns the external IPv6 address for your network.
// Typically, this should be your "temporary" IPv6 address.
func getCurrentIPv6Address(config conf.Config) (string, error) {
	if config.System.IPv6.Enabled == false {
		log.Info(common.MsgIPv6Disabled)
		return "", nil
	}

	log.Println(common.MsgCheckingCurrentIPv6Addr)

	if config.System.IPv6.IPAddrAPI == "" {
		return "", errors.New(common.ErrIPAddressFetchingAPIEmpty)
	}

	resp, err := http.Get(config.System.IPv6.IPAddrAPI)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Body only contains the IP address
	ipv6Address := strings.TrimSpace(string(body))

	log.Println(common.MsgHeaderCurrentIPv6Addr, ipv6Address)

	return ipv6Address, nil
}

func init() {
	flag.StringVar(&conf.ConfigFilePath, "config", "", "Path to the config file.")
	flag.BoolVar(&conf.Debug, "debug", false, "Enable debug logging.")

	flag.Parse()

	log.SetFormatter(&log.TextFormatter{DisableLevelTruncation: true})

	if conf.Debug == true {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
