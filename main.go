package main

import (
	"flag"
	"github.com/boris1993/dnsupdater/conf"
	"github.com/boris1993/dnsupdater/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	var config = conf.Get()

	// Fetch the current external IP address.
	ipAddress := getIPAddr()

	// Then fetch the IP address of the specified DNS record.
	id, recordAddress, err := utils.GetDnsRecordIpAddress()

	if err != nil {
		log.Fatalln(err)
	}

	// Do nothing when the IP address didn't change.
	if ipAddress == recordAddress {
		log.Println("IP address not changed. Will not update the DNS record. ")
		os.Exit(0)
	} else {
		// Update the IP address when changed.
		status, err := utils.UpdateDnsRecord(id, ipAddress)

		if err != nil {
			log.Fatalln(err)
		}

		if !status {
			log.Printf("Failed to update the DNS record %s.\n", config.CloudFlare.DomainName)
			os.Exit(1)
		} else {
			log.Printf("Successfully updated the DNS record %s.\n", config.CloudFlare.DomainName)
		}

		os.Exit(0)
	}
}

func init() {
	absPath, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatalln(err)
	}

	flag.StringVar(&conf.Path, "config", filepath.Join(absPath, "config.yaml"), "Path to the config file.")
	flag.StringVar(&conf.Path, "c", filepath.Join(absPath, "config.yaml"), "Path to the config file.")
	flag.Parse()
}

func getIPAddr() string {
	var config = conf.Get()

	log.Println("Checking current IP address...")

	resp, err := http.Get(config.System.IPAddrAPI)

	if err != nil {
		log.Fatalln(err)
	}

	// Handle errors when closing the HTTP connection
	defer func() {
		err := resp.Body.Close()

		if err != nil {
			log.Println("Error closing the HTTP connection. ", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	// Body only contains the IP address
	ipAddress := string(body)

	log.Printf("Current IP address is: %s.\n", ipAddress)

	return ipAddress
}
