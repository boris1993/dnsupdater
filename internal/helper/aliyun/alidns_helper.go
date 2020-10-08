package aliyun

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/boris1993/dnsupdater/internal/configs"
	"github.com/boris1993/dnsupdater/internal/constants"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// ProcessRecords takes the configuration as well as the current IP address
// then check and update each DNS record in Aliyun DNS
func ProcessRecords(config configs.Config, currentIPAddress string) {
	log.Println(len(config.AliDNSRecords), constants.MsgAliDNSRecordsFoundSuffix)

	for _, aliDNSRecord := range config.AliDNSRecords {
		if aliDNSRecord.RegionID == "" ||
			aliDNSRecord.AccessKeyID == "" ||
			aliDNSRecord.AccessKeySecret == "" ||
			aliDNSRecord.DomainName == "" {
			// Print error and skip to next record when bad configuration found
			log.Errorln(constants.ErrAliDNSRecordConfigIncomplete)
			continue
		}

		// Create an Aliyun API client
		client, err := alidns.NewClientWithAccessKey(aliDNSRecord.RegionID, aliDNSRecord.AccessKeyID, aliDNSRecord.AccessKeySecret)

		if err != nil {
			log.Errorln(err)
			continue
		}

		log.Println(constants.MsgHeaderDomainProcessing, aliDNSRecord.DomainName)

		recordId, recordIP, err := getDomainRecordID(
			aliDNSRecord.DomainName,
			client)

		if err != nil {
			log.Errorln(err)
			continue
		}

		if recordIP == currentIPAddress {
			log.Println(constants.MsgIPAddrNotChanged)
			continue
		}

		// RR is the 2nd level domain
		RR := strings.Split(aliDNSRecord.DomainName, ".")[0]

		status, err := updateAliDNSRecord(recordId, RR, currentIPAddress, client)

		if err != nil {
			log.Errorln(err)
			continue
		}

		if !status {
			log.Errorln(constants.ErrMsgHeaderUpdateDNSRecordFailed, aliDNSRecord.DomainName)
		} else {
			log.Println(constants.MsgHeaderDNSRecordUpdateSuccessful, aliDNSRecord.DomainName)
		}
	}
}

// getDomainRecordID fetches the information of the specified record.
//
// See document: https://help.aliyun.com/document_detail/29776.html?spm=a2c4g.11186623.2.37.1de5425696HU8m#h2-u8FD4u56DEu53C2u65703
func getDomainRecordID(domainName string, client *alidns.Client) (recordId string, recordAddress string, err error) {
	request := alidns.CreateDescribeDomainRecordsRequest()

	// Separate the domain name into several parts.
	// For a 3-level domain name, like "test.example.com",
	// domainNameParts[0] will be "test",
	// domainNameParts[1] will be "example",
	// and domainNameParts[2] will be "com"
	domainNameParts := strings.Split(domainName, ".")

	request.DomainName = strings.Join([]string{domainNameParts[1], domainNameParts[2]}, ".")
	request.RRKeyWord = domainNameParts[0]

	log.Debugln()
	log.Debugln("==========REQUEST FOR GETTING RECORD ID==========")
	log.Debugln("DomainName = " + request.DomainName)
	log.Debugln("RRKeyWord = " + request.RRKeyWord)
	log.Debugln()

	log.Println(constants.MsgHeaderFetchingIPOfDomain, domainName)

	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		return "", "", err
	}

	log.Debugln("Response:\n" + response.GetHttpContentString())

	switch {
	case len(response.DomainRecords.Record) == 0:
		err := errors.New(constants.ErrNoDNSRecordFoundPrefix + domainName)
		return "", "", err
	case len(response.DomainRecords.Record) != 1:
		err := errors.New(constants.ErrMultipleDNSRecordsFound)
		return "", "", err
	}

	record := response.DomainRecords.Record[0]

	log.Printf(constants.MsgFormatAliDNSFetchResult, domainName, record.Value, record.RecordId)

	return record.RecordId, response.DomainRecords.Record[0].Value, nil
}

// Updates the IP address of the specified domain.
//
// See document: https://help.aliyun.com/document_detail/29774.html?spm=a2c4g.11186623.2.35.1de5425696HU8m
func updateAliDNSRecord(recordId string, RR string, currentIPAddress string, client *alidns.Client) (bool, error) {
	request := alidns.CreateUpdateDomainRecordRequest()

	request.RecordId = recordId
	request.RR = RR
	request.Type = "A"
	request.Value = currentIPAddress

	log.Debugln()
	log.Debugln("==========REQUEST FOR UPDATING RECORD ID==========")
	log.Debugln("RecordID = " + request.RecordId)
	log.Debugln("RR = " + request.RR)
	log.Debugln("Type = " + request.Type)
	log.Debugln("Value = " + request.Value)
	log.Debugln()

	log.Printf(constants.MsgFormatUpdatingDNS, recordId, currentIPAddress)

	response, err := client.UpdateDomainRecord(request)

	if err != nil {
		return false, err
	}

	// See document: https://help.aliyun.com/document_detail/29774.html?spm=a2c4g.11186623.2.35.1de5425696HU8m#h2-u9519u8BEFu78014
	httpStatus := response.GetHttpStatus()
	switch httpStatus {
	case http.StatusOK:
		return true, nil
	default:
		return false, errors.New(response.GetHttpContentString())
	}
}
