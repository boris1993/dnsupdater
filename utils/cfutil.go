// Package utils provides utilities about IP addresses, DNS records, and config files.
package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/boris1993/dnsupdater/constants"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/boris1993/dnsupdater/model"
)

// GetDnsRecordIpAddress gets the IP address in the specified DNS record,
// which is identified by the combination of the record type(hard coded as A type for now) and the domain name.
//
// It returns the ID of this DNS record, the IP address of this record,
// or the error message if any error occurs.
func GetDnsRecordIpAddress(conf model.Config) (recordId string, address string, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet,
		conf.CloudFlare.APIEndpoint+"/zones/"+conf.CloudFlare.ZoneID+"/dns_records?type=A&name="+conf.CloudFlare.DomainName,
		nil)

	req.Header.Add("X-Auth-Email", conf.CloudFlare.AuthEmail)
	req.Header.Add("X-Auth-Key", conf.CloudFlare.APIKey)
	req.Header.Add("Content-Type", "application/json")

	log.Printf("Fetching IP address of domain %s.\n", conf.CloudFlare.DomainName)

	resp, err := client.Do(req)

	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != 200 {
		return "", "", errors.New(resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", "", err
	}

	dnsRecord := model.CfDnsRecord{}

	json.Unmarshal([]byte(string(body)), &dnsRecord)

	if !dnsRecord.Success {
		return
	}

	if len(dnsRecord.Result) == 0 {
		return "", "", errors.New(constants.ErrDomainNameNotExist)
	}

	id := dnsRecord.Result[0].ID
	ipAddrInDns := dnsRecord.Result[0].Content

	log.Printf("IP address of %s is %s.\n", conf.CloudFlare.DomainName, ipAddrInDns)

	return id, ipAddrInDns, nil
}

// UpdateDnsRecord updates the specified DNS record identified by the record ID.
//
// id is the record ID, address is the IP address to be written.
//
// It returns the status of the update process, or the error if any error occurs.
func UpdateDnsRecord(id string, address string, conf model.Config) (status bool, err error) {
	client := &http.Client{}

	updateRecordData := model.UpdateRecordData{}
	updateRecordData.RecordType = "A"
	updateRecordData.Name = conf.CloudFlare.DomainName
	updateRecordData.Content = address

	updateRecordDataByte, _ := json.Marshal(updateRecordData)
	requestBodyReader := bytes.NewReader(updateRecordDataByte)

	req, err := http.NewRequest(http.MethodPut,
		conf.CloudFlare.APIEndpoint+"/zones/"+conf.CloudFlare.ZoneID+"/dns_records/"+id,
		requestBodyReader)

	req.Header.Add("X-Auth-Email", conf.CloudFlare.AuthEmail)
	req.Header.Add("X-Auth-Key", conf.CloudFlare.APIKey)
	req.Header.Add("Content-Type", "application/json")

	log.Printf("Updating IP address of domain %s to %s.\n", conf.CloudFlare.DomainName, address)

	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	dnsRecord := model.UpdateRecordResult{}

	json.Unmarshal(body, &dnsRecord)

	updateStatus := dnsRecord.Success

	return updateStatus, nil
}
