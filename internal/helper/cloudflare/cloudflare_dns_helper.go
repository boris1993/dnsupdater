package cloudflare

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/boris1993/dnsupdater/internal/configs"
	"github.com/boris1993/dnsupdater/internal/constants"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// ProcessRecords takes the configuration as well as the current IP address,
// then check and update each DNS record in CloudFlare
func ProcessRecords(currentIPv4Address string, currentIPv6Address string) error {
	config, err := configs.Get()
	if err != nil {
		return err
	}

	if config.System.CloudFlareAPIEndpoint == "" {
		return errors.New(constants.ErrCloudFlareAPIAddressEmpty)
	}

	log.Println(len(config.CloudFlareRecords), constants.MsgCloudFlareRecordsFoundSuffix)

	// Process each CloudFlare DNS record
	for _, cloudFlareRecord := range config.CloudFlareRecords {

		if cloudFlareRecord.APIKey == "" ||
			cloudFlareRecord.AuthEmail == "" ||
			cloudFlareRecord.DomainName == "" ||
			cloudFlareRecord.ZoneID == "" ||
			cloudFlareRecord.DomainType == "" {
			// Print error and skip to next record when bad configuration found
			log.Errorln(constants.ErrCloudFlareRecordConfigIncomplete)
			continue
		}

		if cloudFlareRecord.DomainType != "A" && cloudFlareRecord.DomainType != "AAAA" {
			log.Errorln(constants.ErrInvalidDomainType)
			continue
		}

		// Prints which record is being processed
		log.Println(constants.MsgHeaderDomainProcessing, cloudFlareRecord.DomainName)

		// Then fetch the IP address of the specified DNS record
		id, recordAddress, err := getCFDnsRecordIpAddress(cloudFlareRecord)

		if err != nil {
			log.Errorln(err)
			continue
		}

		// Do nothing when the IP address didn't change.
		if (cloudFlareRecord.DomainType == "A" && currentIPv4Address == recordAddress) ||
			(cloudFlareRecord.DomainType == "AAAA" && currentIPv6Address == recordAddress) {
			log.Println(constants.MsgIPAddrNotChanged)
			continue
		} else {
			var status bool
			var err error

			// Update the IP address when changed.
			switch cloudFlareRecord.DomainType {
			case "A":
				status, err = updateCFDNSRecord(id, currentIPv4Address, cloudFlareRecord)
				break
			case "AAAA":
				status, err = updateCFDNSRecord(id, currentIPv6Address, cloudFlareRecord)
				break
			}

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

// getCFDnsRecordIpAddress gets the IP address associated with the specified DNS record,
// which is identified by the combination of the record type, and the domain name.
//
// cloudFlareRecord contains the information, which are needed by this process, and it is coming from the config.yaml.
//
// The first value returned is the ID of this DNS record,
// the second value returned is the IP address of this record,
// or an error will be returned if any error occurs.
func getCFDnsRecordIpAddress(cloudFlareRecord configs.CloudFlare) (string, string, error) {
	config, err := configs.Get()
	if err != nil {
		return "", "", err
	}

	APIEndpoint := config.System.CloudFlareAPIEndpoint

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet,
		APIEndpoint+"/zones/"+cloudFlareRecord.ZoneID+"/dns_records"+
			"?type="+cloudFlareRecord.DomainType+
			"&name="+cloudFlareRecord.DomainName,
		nil)

	if err != nil {
		return "", "", err
	}

	log.Debug("Request URI: \n" + req.URL.String())

	req.Header.Add("X-Auth-Email", cloudFlareRecord.AuthEmail)
	req.Header.Add("Authorization", "Bearer "+cloudFlareRecord.APIKey)
	req.Header.Add("Content-Type", "application/json")

	log.Println(constants.MsgHeaderFetchingIPOfDomain, cloudFlareRecord.DomainName)

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	cfAPIResponse := CfAPIResponse{}
	err = json.Unmarshal(body, &cfAPIResponse)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != http.StatusOK || !cfAPIResponse.Success {
		errorMessages := ""
		for key := range cfAPIResponse.Errors {
			errorMessages += cfAPIResponse.Errors[key].Message

			if key != len(cfAPIResponse.Errors)-1 {
				errorMessages += "; "
			}
		}

		return "", "", errors.New(errorMessages)
	}

	defer func() {
		err = resp.Body.Close()

		if err != nil {
			log.Errorln(constants.ErrCloseHTTPConnectionFail, err)
		}
	}()

	if len(cfAPIResponse.Result) == 0 {
		return "", "", errors.New(constants.ErrNoDNSRecordFoundPrefix + cloudFlareRecord.DomainName)
	}

	if !cfAPIResponse.Success {
		return "", "", errors.New(constants.ErrMsgHeaderFetchDomainInfoFailed + cloudFlareRecord.DomainName)
	}

	id := cfAPIResponse.Result[0].ID
	ipAddrInDns := cfAPIResponse.Result[0].Content

	log.Printf(constants.MsgFormatDNSFetchResult, cloudFlareRecord.DomainName, ipAddrInDns)

	return id, ipAddrInDns, nil
}

// updateCFDNSRecord updates the specified DNS record identified by the record ID.
//
// id is the record ID, address is the IP address to be written, and cloudFlareRecord contains the information corresponding to the DNS record to be updated.
//
// The return value is the status(true or false) of the update process,
// or an error will be returned if any error occurs.
func updateCFDNSRecord(id string, address string, cloudFlareRecord configs.CloudFlare) (bool, error) {
	config, err := configs.Get()
	if err != nil {
		return false, err
	}

	APIEndpoint := config.System.CloudFlareAPIEndpoint

	client := &http.Client{}

	updateRecordData := UpdateRecordData{}
	updateRecordData.RecordType = cloudFlareRecord.DomainType
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
	req.Header.Add("Authorization", "Bearer "+cloudFlareRecord.APIKey)
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

	dnsRecord := UpdateRecordResult{}

	log.Debug("Response: \n" + string(body))

	err = json.Unmarshal(body, &dnsRecord)

	if err != nil {
		return false, err
	}

	return dnsRecord.Success, nil
}
