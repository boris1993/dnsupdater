// Package cfutil provides utilities about manipulating a CloudFlare DNS record.
package cfutil

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

// ProcessRecords takes the configuration and the current IP address,
// then check and update each DNS record in CloudFlare
func ProcessRecords(config *conf.Config, currentIPAddress string) error {

	if config.System.CloudFlareAPIEndpoint == "" {
		return errors.New(constants.ErrCloudFlareAPIAddressEmpty)
	}

	log.Println(len(config.CloudFlareRecords), constants.MsgCloudFlareRecordsFoundSuffix)

	// Process each CloudFlare DNS record
	for _, cloudFlareRecord := range config.CloudFlareRecords {

		if cloudFlareRecord.APIKey == "" ||
			cloudFlareRecord.AuthEmail == "" ||
			cloudFlareRecord.DomainName == "" ||
			cloudFlareRecord.ZoneID == "" {
			// Print error and skip to next record when bad configuration found
			log.Errorln(constants.ErrCloudFlareRecordConfigIncomplete)
			continue
		}

		// Prints which record is being processed
		log.Println(constants.MsgHeaderDomainProcessing, cloudFlareRecord.DomainName)

		// Then fetch the IP address of the specified DNS record.
		id, recordAddress, err := GetDnsRecordIpAddress(cloudFlareRecord)

		if err != nil {
			log.Errorln(err)
		}

		// Do nothing when the IP address didn't change.
		if currentIPAddress == recordAddress {
			log.Println(constants.MsgIPAddrNotChanged)
			continue
		} else {
			// Update the IP address when changed.
			status, err := UpdateDnsRecord(id, currentIPAddress, cloudFlareRecord)

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

	return nil
}

// GetDnsRecordIpAddress gets the IP address associated with the specified DNS record,
// which is identified by the combination of the record type(hard coded as A type for now) and the domain name.
//
// cloudFlareRecord contains the information which this process needed, and it is coming from the config.yaml.
//
// It returns the ID of this DNS record, the IP address of this record,
// or the error message if any error occurs.
func GetDnsRecordIpAddress(cloudFlareRecord conf.CloudFlare) (recordID string, address string, err error) {
	APIEndpoint := conf.Get().System.CloudFlareAPIEndpoint

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet,
		APIEndpoint+"/zones/"+cloudFlareRecord.ZoneID+"/dns_records?type=A&name="+cloudFlareRecord.DomainName,
		nil)

	req.Header.Add("X-Auth-Email", cloudFlareRecord.AuthEmail)
	req.Header.Add("X-Auth-Key", cloudFlareRecord.APIKey)
	req.Header.Add("Content-Type", "application/json")

	log.Println(constants.MsgHeaderFetchingIPOfDomain, cloudFlareRecord.DomainName)

	resp, err := client.Do(req)

	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != 200 {
		return "", "", errors.New(resp.Status)
	}

	defer func() {
		err = resp.Body.Close()

		if err != nil {
			log.Errorln(constants.ErrCloseHTTPConnectionFail, err)
		}
	}()

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
		return "", "", errors.New(constants.ErrMsgHeaderFetchDomainInfoFailed + cloudFlareRecord.DomainName)
	}

	if len(dnsRecord.Result) == 0 {
		return "", "", errors.New(constants.ErrDomainNameNotExist)
	}

	id := dnsRecord.Result[0].ID
	ipAddrInDns := dnsRecord.Result[0].Content

	log.Printf(constants.MsgFormatDNSFetchResult, cloudFlareRecord.DomainName, ipAddrInDns)

	return id, ipAddrInDns, nil
}

// UpdateDnsRecord updates the specified DNS record identified by the record ID.
//
// id is the record ID, address is the IP address to be written,
// and cloudFlareRecord contains the information corresponding to the DNS record to be updated.
//
// It returns the status of the update process, or the error if any error occurs.
func UpdateDnsRecord(id string, address string, cloudFlareRecord conf.CloudFlare) (status bool, err error) {
	APIEndpoint := conf.Get().System.CloudFlareAPIEndpoint

	client := &http.Client{}

	updateRecordData := model.UpdateRecordData{}
	updateRecordData.RecordType = "A"
	updateRecordData.Name = cloudFlareRecord.DomainName
	updateRecordData.Content = address

	updateRecordDataByte, _ := json.Marshal(updateRecordData)
	requestBodyReader := bytes.NewReader(updateRecordDataByte)

	req, err := http.NewRequest(http.MethodPut,
		APIEndpoint+"/zones/"+cloudFlareRecord.ZoneID+"/dns_records/"+id,
		requestBodyReader)

	if err != nil {
		return false, err
	}

	req.Header.Add("X-Auth-Email", cloudFlareRecord.AuthEmail)
	req.Header.Add("X-Auth-Key", cloudFlareRecord.APIKey)
	req.Header.Add("Content-Type", "application/json")

	log.Printf(constants.MsgFormatUpdatingDNS, cloudFlareRecord.DomainName, address)

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
