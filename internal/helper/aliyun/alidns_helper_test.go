package aliyun

import (
	"encoding/json"
	"github.com/boris1993/dnsupdater/internal/common"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testResourcePath = "../../test"

var testHTTPServer *httptest.Server
var config *common.Config

var serverRecords []AliDNSRecord

func TestAliDNSHelper(t *testing.T) {
	var err error

	err = prepareMockedServerResponse()
	if err != nil {
		t.Error(err)
		return
	}

	err = startTestHTTPServer()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("testProcessRecords", testProcessRecords)

	stopTestHTTPServer()
}

func testProcessRecords(t *testing.T) {
	currentIPAddress := "192.168.1.1"
	currentIPv6Address := "240a:38b:5dc0:5d00:e1dd:c7c7:169a:acbb"

	err := ProcessRecords(currentIPAddress, currentIPv6Address)
	if err != nil {
		t.Error(err)
		return
	}
}

func prepareMockedServerResponse() error {
	mockServerResponseFilePath := testResourcePath + "/mock_aliyun_dns_response.json"
	bytes, err := ioutil.ReadFile(mockServerResponseFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &serverRecords)
	if err != nil {
		return err
	}

	return nil
}

func startTestHTTPServer() error {
	httpHandler := http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			mockAliyunDNSResponse(writer, request)
		},
	)

	testHTTPServer = httptest.NewServer(httpHandler)

	err := setEndpointToTestServer()
	if err != nil {
		return err
	}

	return nil
}

func stopTestHTTPServer() {
	testHTTPServer.Close()
}

func setEndpointToTestServer() error {
	var err error
	common.ConfigFilePath = testResourcePath + "/test_config.yaml"

	config, err = common.GetConfig()
	if err != nil {
		return err
	}
	config.System.AliyunAPIEndpoint = testHTTPServer.URL

	return nil
}

func mockAliyunDNSResponse(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Query().Get("Action") {
	case ActionDescribeDomainRecords:
		domainName := request.URL.Query().Get(QueryParamDomainName)
		keyWord := request.URL.Query().Get(QueryParamKeyWord)
		domainType := request.URL.Query().Get(QueryParamType)

		for index := range serverRecords {
			serverRecord := serverRecords[index]
			if serverRecord.DomainName == domainName &&
				serverRecord.RR == keyWord &&
				serverRecord.Type == domainType {
				_, err := writer.Write(generateSuccessDescribeDomainRecordResponseJson(serverRecord))
				if err != nil {
					log.Error(err)
					return
				}
				return
			}
		}

		_, err := writer.Write(generateEmptyDescribeDomainRecordResponseJson())
		if err != nil {
			log.Error(err)
			return
		}
		return
	case ActionUpdateDomainRecord:
		recordId := request.URL.Query().Get(QueryParamRecordId)

		for index := range serverRecords {
			if serverRecords[index].RecordID == recordId {
				_, err := writer.Write([]byte{})
				if err != nil {
					log.Error(err)
					return
				}
				break
			}
		}
	}
}

func generateSuccessDescribeDomainRecordResponseJson(serverRecord AliDNSRecord) []byte {
	describeDomainRecordsResponse := DescribeDomainRecordsResponse{
		commonResponse: commonResponse{RequestID: "dummyRequestID"},
		TotalCount:     1,
		PageNumber:     1,
		PageSize:       20,
		DomainRecords: domainRecords{
			[]AliDNSRecord{serverRecord},
		},
	}

	jsonByte, _ := json.Marshal(describeDomainRecordsResponse)

	return jsonByte
}

func generateEmptyDescribeDomainRecordResponseJson() []byte {
	describeDomainRecordsResponse := DescribeDomainRecordsResponse{
		commonResponse: commonResponse{RequestID: "dummyRequestID"},
		TotalCount:     0,
		PageNumber:     1,
		PageSize:       20,
		DomainRecords: domainRecords{
			[]AliDNSRecord{},
		},
	}

	jsonByte, _ := json.Marshal(describeDomainRecordsResponse)

	return jsonByte
}
