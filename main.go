package main

import (
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
	var config = conf.Get()

	// Fetch the current external IP address.
	ipAddress := getIPAddr()

	log.Println(len(config.CloudFlareRecords), constants.MsgCloudFlareRecordsFoundSuffix)

	// Process each CloudFlare DNS record
	for _, cloudFlareRecord := range config.CloudFlareRecords {

		// Prints which record is being processed
		log.Println(constants.MsgHeaderDomainProcessing, cloudFlareRecord.DomainName)

		// Then fetch the IP address of the specified DNS record.
		id, recordAddress, err := cfutil.GetDnsRecordIpAddress(cloudFlareRecord)

		if err != nil {
			log.Errorln(err)
		}

		// Do nothing when the IP address didn't change.
		if ipAddress == recordAddress {
			log.Println(constants.MsgIPAddrNotChanged)
			continue
		} else {
			// Update the IP address when changed.
			status, err := cfutil.UpdateDnsRecord(id, ipAddress, cloudFlareRecord)

			if err != nil {
				log.Errorln(err)
				continue
			}

			if !status {
				log.Errorln(constants.ErrMsgHeaderUpdateDNSRecordFailed, cloudFlareRecord.DomainName)
				continue
			} else {
				log.Println(constants.MsgHeaderDNSRecordUpdateSuccessful, cloudFlareRecord.DomainName)
			}
		}
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

// getIPAddr returns the external IP address for your network
func getIPAddr() string {
	var config = conf.Get()

	log.Println(constants.MsgCheckingCurrentIPAddr)

	resp, err := http.Get(config.System.IPAddrAPI)

	if err != nil {
		log.Fatalln(err)
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
		log.Fatalln(err)
	}

	// Body only contains the IP address
	ipAddress := string(body)

	log.Println(constants.MsgHeaderCurrentIPAddr, ipAddress)

	return ipAddress
}
