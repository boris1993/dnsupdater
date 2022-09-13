package aliyun

import (
	"encoding/json"
	"errors"
	"github.com/boris1993/dnsupdater/internal/common"
	"github.com/boris1993/dnsupdater/internal/conf"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

// ProcessRecords takes the configuration as well as the current IP address
// then check and update each DNS record in Aliyun DNS
func ProcessRecords(currentIPv4Address string, currentIPv6Address string) error {
	config, err := conf.GetConfig()
	if err != nil {
		return err
	}

	if config.System.AliyunAPIEndpoint == "" {
		return errors.New(common.ErrAliyunAPIAddressEmpty)
	}

	log.Println(len(config.AliDNSRecords), common.MsgAliDNSRecordsFoundSuffix)

	for _, aliDNSRecord := range config.AliDNSRecords {
		if aliDNSRecord.RegionID == "" ||
			aliDNSRecord.AccessKeyID == "" ||
			aliDNSRecord.AccessKeySecret == "" ||
			aliDNSRecord.DomainName == "" ||
			aliDNSRecord.DomainType == "" {
			// Print error and skip to next record when bad configuration found
			log.Errorln(common.ErrAliDNSRecordConfigIncomplete)
			continue
		}

		if aliDNSRecord.DomainType != "A" && aliDNSRecord.DomainType != "AAAA" {
			log.Errorln(common.ErrInvalidDomainType)
			continue
		}

		log.Println(common.MsgHeaderDomainProcessing, aliDNSRecord.DomainName)

		recordId, recordIP, err := getDomainRecordID(aliDNSRecord)
		if err != nil {
			log.Errorln(err)
			continue
		}

		if (aliDNSRecord.DomainType == "A" && common.CompareAddresses(currentIPv4Address, recordIP)) ||
			(aliDNSRecord.DomainType == "AAAA" && common.CompareAddresses(currentIPv6Address, recordIP)) {
			log.Println(common.MsgIPAddrNotChanged)
			continue
		}

		domainNameParts := strings.Split(aliDNSRecord.DomainName, ".")
		// Remove the TLD and the domain, what's rest are the host record
		hostRecord := strings.Join(domainNameParts[:len(domainNameParts)-2], ".")

		var status bool
		switch aliDNSRecord.DomainType {
		case "A":
			status, err = updateAliDNSRecord(recordId, hostRecord, currentIPv4Address, aliDNSRecord)
			break
		case "AAAA":
			// If there's no valid IPv6 internet address,
			// then skip updating this record and head to the next one
			if currentIPv6Address == "" {
				log.Info(common.MsgIPv6AddrNotAvailable)
				continue
			}

			status, err = updateAliDNSRecord(recordId, hostRecord, currentIPv6Address, aliDNSRecord)
			break
		}

		if err != nil {
			log.Errorln(err)
			continue
		}

		if !status {
			log.Errorln(common.ErrMsgHeaderUpdateDNSRecordFailed, aliDNSRecord.DomainName)
		} else {
			log.Println(common.MsgHeaderDNSRecordUpdateSuccessful, aliDNSRecord.DomainName)
		}
	}

	return nil
}

// getDomainRecordID fetches the information of the specified record.
//
// See document: https://help.aliyun.com/document_detail/29776.html?spm=a2c4g.11186623.2.37.1de5425696HU8m#h2-u8FD4u56DEu53C2u65703
func getDomainRecordID(aliDNSConfigRecord conf.AliDNS) (recordId string, recordAddress string, err error) {
	config, err := conf.GetConfig()
	if err != nil {
		return "", "", err
	}

	domainNameParts := strings.Split(aliDNSConfigRecord.DomainName, ".")

	tld := domainNameParts[len(domainNameParts)-1]
	domainName := domainNameParts[len(domainNameParts)-2]
	hostRecord := strings.Join(domainNameParts[:len(domainNameParts)-2], ".")

	param := make(map[string]string)
	param[QueryParamDomainName] = strings.Join([]string{domainName, tld}, ".")
	param[QueryParamKeyWord] = hostRecord
	param[QueryParamSearchMode] = SearchModeAdvanced
	param[QueryParamType] = aliDNSConfigRecord.DomainType

	request, err := BuildAliDNSRequest(
		config.System.AliyunAPIEndpoint,
		aliDNSConfigRecord.AccessKeyID,
		aliDNSConfigRecord.AccessKeySecret,
		ActionDescribeDomainRecords,
		param,
	)

	log.Println(common.MsgHeaderFetchingIPOfDomain, aliDNSConfigRecord.DomainName)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return "", "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer func() {
		err := response.Body.Close()

		if err != nil {
			log.Errorln(common.ErrCloseHTTPConnectionFail, err)
		}
	}()
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

	describeDomainRecordsResponse := DescribeDomainRecordsResponse{}
	err = json.Unmarshal(body, &describeDomainRecordsResponse)
	if err != nil {
		return "", "", err
	}

	if len(describeDomainRecordsResponse.DomainRecords.Record) == 0 {
		err := errors.New(common.ErrNoDNSRecordFoundPrefix + aliDNSConfigRecord.DomainName)
		return "", "", err
	}

	// The API is poorly designed and I just can't figure out a way to exactly retrieve the wanted record,
	// so I have to iterate the response and find the matched one.
	for _, domainRecord := range describeDomainRecordsResponse.DomainRecords.Record {
		if domainRecord.RR == hostRecord && domainRecord.Type == aliDNSConfigRecord.DomainType {
			log.Printf(common.MsgFormatAliDNSFetchResult,
				aliDNSConfigRecord.DomainName,
				domainRecord.Value,
				domainRecord.RecordID)

			return domainRecord.RecordID, domainRecord.Value, nil
		}
	}

	return "",
		"",
		errors.New(common.ErrMsgHeaderFetchDomainInfoFailed + aliDNSConfigRecord.DomainName)
}

// Updates the IP address of the specified domain.
//
// See document: https://help.aliyun.com/document_detail/29774.html?spm=a2c4g.11186623.2.35.1de5425696HU8m
func updateAliDNSRecord(recordId string, RR string, currentIPAddress string, aliDNSConfigRecord conf.AliDNS) (bool, error) {
	config, err := conf.GetConfig()
	if err != nil {
		return false, err
	}

	param := make(map[string]string)
	param[QueryParamRecordId] = recordId
	param[QueryParamRR] = RR
	param[QueryParamType] = aliDNSConfigRecord.DomainType
	param[QueryParamValue] = currentIPAddress

	request, err := BuildAliDNSRequest(
		config.System.AliyunAPIEndpoint,
		aliDNSConfigRecord.AccessKeyID,
		aliDNSConfigRecord.AccessKeySecret,
		ActionUpdateDomainRecord,
		param,
	)

	log.Printf(common.MsgFormatUpdatingDNS, recordId, currentIPAddress)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return false, err
	}

	// See document: https://help.aliyun.com/document_detail/29774.html?spm=a2c4g.11186623.2.35.1de5425696HU8m#h2-u9519u8BEFu78014
	httpStatus := response.StatusCode
	switch httpStatus {
	case http.StatusOK:
		return true, nil
	default:
		body, err := ioutil.ReadAll(response.Body)
		defer func() {
			err := response.Body.Close()

			if err != nil {
				log.Errorln(common.ErrCloseHTTPConnectionFail, err)
			}
		}()
		if err != nil {
			return false, err
		}

		aliDNSErrorResponse := AliDNSErrorResponse{}
		err = json.Unmarshal(body, &aliDNSErrorResponse)
		if err != nil {
			return false, err
		}

		return false, errors.New(aliDNSErrorResponse.Message)
	}
}
