package aliyun

import (
	"encoding/json"
	"errors"
	"github.com/boris1993/dnsupdater/internal/configs"
	"github.com/boris1993/dnsupdater/internal/constants"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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

		log.Println(constants.MsgHeaderDomainProcessing, aliDNSRecord.DomainName)

		recordId, recordIP, err := getDomainRecordID(aliDNSRecord)

		if err != nil {
			log.Errorln(err)
			continue
		}

		if recordIP == currentIPAddress {
			log.Println(constants.MsgIPAddrNotChanged)
			continue
		}

		// RR is the 2nd level domain
		domainNameParts := strings.Split(aliDNSRecord.DomainName, ".")
		RR := strings.Join(domainNameParts[:len(domainNameParts)-2], ".")

		status, err := updateAliDNSRecord(recordId, RR, currentIPAddress, aliDNSRecord)

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
func getDomainRecordID(aliDNSConfigRecord configs.AliDNS) (recordId string, recordAddress string, err error) {
	domainNameParts := strings.Split(aliDNSConfigRecord.DomainName, ".")

	tld := domainNameParts[len(domainNameParts)-1]
	domainName := domainNameParts[len(domainNameParts)-2]
	hostRecord := strings.Join(domainNameParts[:len(domainNameParts)-2], ".")

	param := make(map[string]string)
	param["DomainName"] = strings.Join([]string{domainName, tld}, ".")
	param["KeyWord"] = hostRecord
	param["SearchMode"] = "EXACT"

	request, err := BuildAliDNSRequest(
		configs.Get().System.AliyunAPIEndpoint,
		aliDNSConfigRecord.AccessKeyID,
		aliDNSConfigRecord.AccessKeySecret,
		"DescribeDomainRecords",
		param,
	)

	if log.GetLevel() == log.DebugLevel {
		log.Debugln()
		log.Debugln("==========REQUEST FOR GETTING RECORD ID==========")
		for key, value := range param {
			log.Debugf("%s = %s", key, value)
		}
		log.Debugln("==========REQUEST FOR GETTING RECORD ID==========")
		log.Debugln()
	}

	log.Println(constants.MsgHeaderFetchingIPOfDomain, aliDNSConfigRecord.DomainName)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return "", "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	if response.StatusCode != http.StatusOK {
		aliDNSErrorResponse := AliDNSErrorResponse{}
		err = json.Unmarshal(body, &aliDNSErrorResponse)
		if err != nil {
			return "", "", err
		}

		return "", "", errors.New(aliDNSErrorResponse.Message)
	}

	defer func() {
		err := response.Body.Close()

		if err != nil {
			log.Errorln(constants.ErrCloseHTTPConnectionFail, err)
		}
	}()

	describeDomainRecordsResponse := DescribeDomainRecordsResponse{}
	err = json.Unmarshal(body, &describeDomainRecordsResponse)
	if err != nil {
		return "", "", err
	}

	switch {
	case len(describeDomainRecordsResponse.DomainRecords.Record) == 0:
		err := errors.New(constants.ErrNoDNSRecordFoundPrefix + aliDNSConfigRecord.DomainName)
		return "", "", err
	case len(describeDomainRecordsResponse.DomainRecords.Record) != 1:
		err := errors.New(constants.ErrMultipleDNSRecordsFound)
		return "", "", err
	}

	record := describeDomainRecordsResponse.DomainRecords.Record[0]

	log.Printf(constants.MsgFormatAliDNSFetchResult, record.DomainName, record.Value, record.RecordID)

	return record.RecordID, record.Value, nil
}

// Updates the IP address of the specified domain.
//
// See document: https://help.aliyun.com/document_detail/29774.html?spm=a2c4g.11186623.2.35.1de5425696HU8m
func updateAliDNSRecord(recordId string, RR string, currentIPAddress string, aliDNSConfigRecord configs.AliDNS) (bool, error) {
	param := make(map[string]string)
	param["RecordId"] = recordId
	param["RR"] = RR
	param["Type"] = "A"
	param["Value"] = currentIPAddress

	request, err := BuildAliDNSRequest(
		configs.Get().System.AliyunAPIEndpoint,
		aliDNSConfigRecord.AccessKeyID,
		aliDNSConfigRecord.AccessKeySecret,
		"UpdateDomainRecord",
		param,
	)

	if log.GetLevel() == log.DebugLevel {
		log.Debugln()
		log.Debugln("==========REQUEST FOR UPDATING RECORD ID==========")
		for key, value := range param {
			log.Debugf("%s = %s", key, value)
		}
		log.Debugln("==========REQUEST FOR UPDATING RECORD ID==========")
		log.Debugln()
	}

	log.Printf(constants.MsgFormatUpdatingDNS, recordId, currentIPAddress)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	// See document: https://help.aliyun.com/document_detail/29774.html?spm=a2c4g.11186623.2.35.1de5425696HU8m#h2-u9519u8BEFu78014
	httpStatus := response.StatusCode
	switch httpStatus {
	case http.StatusOK:
		return true, nil
	default:
		aliDNSErrorResponse := AliDNSErrorResponse{}
		err = json.Unmarshal(body, &aliDNSErrorResponse)
		if err != nil {
			return false, err
		}
		return false, errors.New(aliDNSErrorResponse.Message)
	}
}
