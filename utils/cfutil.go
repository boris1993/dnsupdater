// Package utils provides utilities about IP addresses, DNS records, and config files.
package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/boris1993/dnsupdater/conf"
	"github.com/boris1993/dnsupdater/constants"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"

	"github.com/boris1993/dnsupdater/model"
)

// GetDnsRecordIpAddress gets the IP address in the specified DNS record,
// which is identified by the combination of the record type(hard coded as A type for now) and the domain name.
//
// It returns the ID of this DNS record, the IP address of this record,
// or the error message if any error occurs.
func GetDnsRecordIpAddress() (recordID string, address string, err error) {
	var config = conf.Get()

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet,
		config.CloudFlare.APIEndpoint+"/zones/"+config.CloudFlare.ZoneID+"/dns_records?type=A&name="+config.CloudFlare.DomainName,
		nil)

	req.Header.Add("X-Auth-Email", config.CloudFlare.AuthEmail)
	req.Header.Add("X-Auth-Key", config.CloudFlare.APIKey)
	req.Header.Add("Content-Type", "application/json")

	log.Println(constants.MsgHeaderFetchingIPOfDomain, config.CloudFlare.DomainName)

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

	log.Debug("Response: \n" + string(body))

	err = json.Unmarshal([]byte(string(body)), &dnsRecord)

	if err != nil {
		return "", "", err
	}

	if !dnsRecord.Success {
		return "", "", errors.New(constants.ErrMsgHeaderFetchDomainInfoFailed + config.CloudFlare.DomainName)
	}

	if len(dnsRecord.Result) == 0 {
		return "", "", errors.New(constants.ErrDomainNameNotExist)
	}

	id := dnsRecord.Result[0].ID
	ipAddrInDns := dnsRecord.Result[0].Content

	log.Printf(constants.MsgFormatDNSFetchResult, config.CloudFlare.DomainName, ipAddrInDns)

	return id, ipAddrInDns, nil
}

// UpdateDnsRecord updates the specified DNS record identified by the record ID.
//
// id is the record ID, address is the IP address to be written.
//
// It returns the status of the update process, or the error if any error occurs.
func UpdateDnsRecord(id string, address string) (status bool, err error) {
	var config = conf.Get()

	client := &http.Client{}

	updateRecordData := model.UpdateRecordData{}
	updateRecordData.RecordType = "A"
	updateRecordData.Name = config.CloudFlare.DomainName
	updateRecordData.Content = address

	updateRecordDataByte, _ := json.Marshal(updateRecordData)
	requestBodyReader := bytes.NewReader(updateRecordDataByte)

	req, err := http.NewRequest(http.MethodPut,
		config.CloudFlare.APIEndpoint+"/zones/"+config.CloudFlare.ZoneID+"/dns_records/"+id,
		requestBodyReader)

	if err != nil {
		return false, err
	}

	req.Header.Add("X-Auth-Email", config.CloudFlare.AuthEmail)
	req.Header.Add("X-Auth-Key", config.CloudFlare.APIKey)
	req.Header.Add("Content-Type", "application/json")

	log.Printf(constants.MsgFormatUpdatingDNS, config.CloudFlare.DomainName, address)

	resp, err := client.Do(req)

	if err != nil {
		return false, err
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
		return false, err
	}

	dnsRecord := model.UpdateRecordResult{}

	log.Debug("Response: \n" + string(body))

	err = json.Unmarshal(body, &dnsRecord)

	if err != nil {
		return false, err
	}

	return dnsRecord.Success, nil
}
