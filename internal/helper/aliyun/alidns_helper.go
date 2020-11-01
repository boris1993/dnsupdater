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
func ProcessRecords(currentIPAddress string) error {
	config, err := configs.Get()
	if err != nil {
		return err
	}

	if config.System.AliyunAPIEndpoint == "" {
		return errors.New(constants.ErrAliyunAPIAddressEmpty)
	}

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

		domainNameParts := strings.Split(aliDNSRecord.DomainName, ".")
		// Remove the TLD and the domain, what's rest are the host record
		hostRecord := strings.Join(domainNameParts[:len(domainNameParts)-2], ".")

		status, err := updateAliDNSRecord(recordId, hostRecord, currentIPAddress, aliDNSRecord)
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

	return nil
}

// getDomainRecordID fetches the information of the specified record.
//
// See document: https://help.aliyun.com/document_detail/29776.html?spm=a2c4g.11186623.2.37.1de5425696HU8m#h2-u8FD4u56DEu53C2u65703
func getDomainRecordID(aliDNSConfigRecord configs.AliDNS) (recordId string, recordAddress string, err error) {
	config, err := configs.Get()
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
	param[QueryParamSearchMode] = SearchModeExact

	request, err := BuildAliDNSRequest(
		config.System.AliyunAPIEndpoint,
		aliDNSConfigRecord.AccessKeyID,
		aliDNSConfigRecord.AccessKeySecret,
		ActionDescribeDomainRecords,
		param,
	)

	log.Println(constants.MsgHeaderFetchingIPOfDomain, aliDNSConfigRecord.DomainName)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return "", "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer func() {
		err := response.Body.Close()

		if err != nil {
			log.Errorln(constants.ErrCloseHTTPConnectionFail, err)
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
		err := errors.New(constants.ErrNoDNSRecordFoundPrefix + aliDNSConfigRecord.DomainName)
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
	config, err := configs.Get()
	if err != nil {
		return false, err
	}

	param := make(map[string]string)
	param[QueryParamRecordId] = recordId
	param[QueryParamRR] = RR
	param[QueryParamType] = "A"
	param[QueryParamValue] = currentIPAddress

	request, err := BuildAliDNSRequest(
		config.System.AliyunAPIEndpoint,
		aliDNSConfigRecord.AccessKeyID,
		aliDNSConfigRecord.AccessKeySecret,
		ActionUpdateDomainRecord,
		param,
	)

	log.Printf(constants.MsgFormatUpdatingDNS, recordId, currentIPAddress)

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
				log.Errorln(constants.ErrCloseHTTPConnectionFail, err)
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
