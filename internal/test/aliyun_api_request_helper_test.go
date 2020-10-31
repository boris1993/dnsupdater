package test

import (
	"github.com/boris1993/dnsupdater/internal/helper/aliyun"
	"net/http"
	"testing"
)

func TestSignRequest(t *testing.T) {
	expectedSignature := "uRpHwaSEt3J+6KQD//svCh/x+pI="

	accessKeySecret := "testsecret"

	queryParams := make(map[string][]string)
	queryParams["Format"] = []string{"XML"}
	queryParams["AccessKeyId"] = []string{"testid"}
	queryParams["Action"] = []string{"DescribeDomainRecords"}
	queryParams["SignatureMethod"] = []string{"HMAC-SHA1"}
	queryParams["DomainName"] = []string{"example.com"}
	queryParams["SignatureNonce"] = []string{"f59ed6a9-83fc-473b-9cc6-99c95df3856e"}
	queryParams["SignatureVersion"] = []string{"1.0"}
	queryParams["Version"] = []string{"2015-01-09"}
	queryParams["Timestamp"] = []string{"2016-03-24T16:41:54Z"}

	signature := aliyun.GenerateAliDNSRequestSignature(http.MethodGet, accessKeySecret, queryParams)

	if signature != expectedSignature {
		t.Errorf("GenerateAliDNSRequestSignature expected = %s, got = %s", expectedSignature, signature)
	}
}

func TestBuildRequest(t *testing.T) {
	apiEndpoint := "https://alidns.aliyuncs.com"
	action := "DescribeDomainRecords"
	accessKeyId := "testid"
	accessKeySecret := "testsecret"

	queryParams := make(map[string]string)
	queryParams["DomainName"] = "example.com"

	_, err := aliyun.BuildAliDNSRequest(apiEndpoint, accessKeyId, accessKeySecret, action, queryParams)

	if err != nil {
		t.Errorf("Error occurred when testing BuildAliDNSRequest. Message=%s", err)
	}
}
