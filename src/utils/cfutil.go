package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"config"
	"model"
)

func GetDnsRecordIpAddress() (recordId string, address string, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet,
		config.CfAPIEndpoint+"/zones/"+config.CfZoneID+"/dns_records?type=A&name="+config.CfDomainName,
		nil)

	req.Header.Add("X-Auth-Email", config.CfAuthEmail)
	req.Header.Add("X-Auth-Key", config.CfAPIKey)
	req.Header.Add("Content-Type", "application/json")

	log.Println("Fetching IP address of domain " + config.CfDomainName)

	resp, err := client.Do(req)

	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", "", err
	}

	dnsRecord := model.CfDnsRecord{}

	json.Unmarshal([]byte(string(body)), &dnsRecord)

	id := dnsRecord.Result[0].Id
	ipAddrInDns := dnsRecord.Result[0].Content

	log.Println("IP address of " + config.CfDomainName + " is " + ipAddrInDns)

	return id, ipAddrInDns, nil
}

func UpdateDnsRecord(id string, address string) (status bool, err error) {
	client := &http.Client{}

	updateRecordData := model.UpdateRecordData{}
	updateRecordData.RecordType = "A"
	updateRecordData.Name = config.CfDomainName
	updateRecordData.Content = address

	updateRecordDataByte, _ := json.Marshal(updateRecordData)
	requestBodyReader := bytes.NewReader(updateRecordDataByte)

	req, err := http.NewRequest(http.MethodPut,
		config.CfAPIEndpoint+"/zones/"+config.CfZoneID+"/dns_records/"+id,
		requestBodyReader)

	req.Header.Add("X-Auth-Email", config.CfAuthEmail)
	req.Header.Add("X-Auth-Key", config.CfAPIKey)
	req.Header.Add("Content-Type", "application/json")

	log.Println("Updating IP address of domain " + config.CfDomainName + " to " + address)

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
