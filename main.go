package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"dnsupdater/config"
	"dnsupdater/utils"
)

func main() {
	ipAddress := getIpAddr()

	id, recordAddress, err := utils.GetDnsRecordIpAddress()

	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	if ipAddress == recordAddress {
		log.Println("IP address not changed. Will not update the DNS record. ")
		os.Exit(0)
	} else {
		log.Println("New IP address is " + ipAddress + ". Updating the DNS record. ")
		status, err := utils.UpdateDnsRecord(id, ipAddress)

		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}

		if !status {
			log.Println("Failed to update the DNS record " + config.CfDomainName + ". Please double check on the CloudFlare dashboard. ")
			os.Exit(-1)
		} else {
			log.Println("Successfully updated the DNS record " + config.CfDomainName + ". ")
		}

		os.Exit(0)
	}
}

func getIpAddr() string {
	log.Println("Checking current IP address...")

	resp, err := http.Get(config.IPAddrAPI)

	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	ipAddress := string(body)

	log.Println("Current IP address is: " + ipAddress)

	return ipAddress
}
