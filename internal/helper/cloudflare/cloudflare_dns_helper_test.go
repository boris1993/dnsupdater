package cloudflare

import (
	"encoding/json"
	"github.com/boris1993/dnsupdater/internal/configs"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testResourcePath = "../../test"

var testHTTPServer *httptest.Server
var config *configs.Config

var serverRecords []cfDnsRecordResult

func TestCloudFlareDNSHelper(t *testing.T) {
	var err error

	err = prepareMockedQueryDomainResponse()
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
	currentIPv4Address := "192.168.1.1"
	currentIPv6Address := "240a:38b:5dc0:5d00:e1dd:c7c7:169a:acbb"

	err := ProcessRecords(currentIPv4Address, currentIPv6Address)
	if err != nil {
		t.Error(err)
		return
	}
}

func prepareMockedQueryDomainResponse() error {
	mockServerResponseFilePath := testResourcePath + "/mock_cloudflare_dns_response.json"
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

func generateMockedUpdateDomainResponse() ([]byte, error) {
	mockServerResponseFilePath := testResourcePath + "/mock_cloudflare_dns_update_response.json"
	bytes, err := ioutil.ReadFile(mockServerResponseFilePath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func startTestHTTPServer() error {
	httpHandler := http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			mockCloudFlareDNSResponse(writer, request)
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
	configs.Path = testResourcePath + "/test_config.yaml"

	config, err = configs.Get()
	if err != nil {
		return err
	}
	config.System.CloudFlareAPIEndpoint = testHTTPServer.URL

	return nil
}

func mockCloudFlareDNSResponse(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	// Handle query domain request
	case http.MethodGet:
		domainName := request.URL.Query().Get("name")
		domainType := request.URL.Query().Get("type")

		for index := range serverRecords {
			if serverRecords[index].Name == domainName &&
				serverRecords[index].RecordType == domainType {
				_, err := writer.Write(generateQueryDomainResponse(serverRecords[index]))
				if err != nil {
					log.Error(err)
					return
				}
				return
			}
		}

	case http.MethodPut:
		// I'm not gonna mock another response for the IPv6 record
		// because we only care about the value of "success"
		bytes, err := generateMockedUpdateDomainResponse()
		if err != nil {
			log.Error(err)
			return
		}

		_, err = writer.Write(bytes)
		if err != nil {
			log.Error(err)
			return
		}
	}

}

func generateQueryDomainResponse(cfDnsRecord cfDnsRecordResult) []byte {
	cfAPIResponse := CfAPIResponse{
		Success:  true,
		Errors:   []errorMessage{},
		Messages: []string{},
		Result:   []cfDnsRecordResult{cfDnsRecord},
		ResultInfo: cfDnsRecordResultInfo{
			Count:      1,
			Page:       1,
			PerPage:    20,
			TotalCount: 1,
			TotalPages: 1,
		},
	}

	jsonByte, _ := json.Marshal(cfAPIResponse)

	return jsonByte
}
