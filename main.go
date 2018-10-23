package main

import (
	"flag"
	"github.com/boris1993/dnsupdater/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/boris1993/dnsupdater/utils"
)

var ConfigPath string

func main() {
	log.Printf("Loading configuraton from %s.\n", ConfigPath)

	// Try to parse the config file first.
	conf, err := utils.ReadConfig(ConfigPath)

	// Print the error message and exit when failed to parse the config file.
	if err != nil {
		log.Fatalln(err)
	}

	// Fetch the current external IP address.
	ipAddress := getIPAddr(conf)

	// Then fetch the IP address of the specified DNS record.
	id, recordAddress, err := utils.GetDnsRecordIpAddress(conf)

	if err != nil {
		log.Fatal(err)
	}

	// Do nothing when the IP address didn't change.
	if ipAddress == recordAddress {
		log.Println("IP address not changed. Will not update the DNS record. ")
		os.Exit(0)
	} else {
		// Update the IP address when changed.
		status, err := utils.UpdateDnsRecord(id, ipAddress, conf)

		if err != nil {
			log.Fatal(err)
		}

		if !status {
			log.Printf("Failed to update the DNS record %s.\n", conf.CloudFlare.DomainName)
			os.Exit(1)
		} else {
			log.Printf("Successfully updated the DNS record %s.\n", conf.CloudFlare.DomainName)
		}

		os.Exit(0)
	}
}

func init() {
	flag.StringVar(&ConfigPath, "config", "./config.yaml", "Path to the config file.")
	flag.StringVar(&ConfigPath, "c", "./config.yaml", "Path to the config file.")
	flag.Parse()
}

func getIPAddr(conf model.Config) string {
	log.Println("Checking current IP address...")

	resp, err := http.Get(conf.System.IPAddrAPI)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	// Body only contains the IP address
	ipAddress := string(body)

	log.Printf("Current IP address is: %s.\n", ipAddress)

	return ipAddress
}
