package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func BuildAliDNSRequest(
	apiEndpoint string,
	accessKeyId string,
	accessKeySecret string,
	action string,
	requestSpecificQueryParameters ...map[string]string,
) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, apiEndpoint, nil)

	if err != nil {
		return nil, err
	}

	queryParams := request.URL.Query()
	//region Common parameters
	queryParams.Add("Action", action)
	queryParams.Add("Format", "JSON")
	queryParams.Add("Version", "2015-01-09")
	queryParams.Add("AccessKeyId", accessKeyId)
	queryParams.Add("SignatureMethod", "HMAC-SHA1")
	queryParams.Add("Timestamp", time.Now().UTC().Format(time.RFC3339))
	queryParams.Add("SignatureVersion", "1.0")
	queryParams.Add("SignatureNonce", uuid.New().String())
	//endregion
	//region Request specific parameters
	for index := range requestSpecificQueryParameters {
		for key, value := range requestSpecificQueryParameters[index] {
			queryParams.Add(key, value)
		}
	}
	//endregion

	signature := GenerateAliDNSRequestSignature(request.Method, accessKeySecret, queryParams)
	queryParams.Add("Signature", signature)

	request.URL.RawQuery = queryParams.Encode()
	return request, nil
}

//GenerateAliDNSRequestSignature returns the signature of the request.
//
//See https://help.aliyun.com/document_detail/29747.html?spm=a2c4g.11186623.6.632.53b98dcazC573U for more information.
func GenerateAliDNSRequestSignature(httpMethod string, accessSecret string, queryParams map[string][]string) string {
	// The keys in the query parameter string must be sorted alphabetically and is case-sensitive.
	keys := make([]string, 0, len(queryParams))
	for key := range queryParams {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Concatenate the sorted keys and its value into a new query parameter string.
	sortedQueryString := ""
	for index := range keys {
		key := keys[index]
		// Each query parameter contains only one value so I'll just pick the first one.
		value := queryParams[key][0]

		// The key and the value must be escaped
		escapedKey := escapePlusAndStarAndUnescapeTildeSymbol(url.QueryEscape(key))
		escapedValue := escapePlusAndStarAndUnescapeTildeSymbol(url.QueryEscape(value))

		sortedQueryString += escapedKey + "=" + escapedValue

		if index != (len(keys) - 1) {
			sortedQueryString += "&"
		}
	}

	// The query parameter string is escaped again as a whole,
	// prefixed with the HTTP method and a string "&%2F&",
	// then we have a new string, and we'll generate its hash as its signature.
	stringToSign := httpMethod + "&%2F&" + url.QueryEscape(sortedQueryString)

	hmacKey := accessSecret + "&"
	hmacHash := hmac.New(sha1.New, []byte(hmacKey))
	hmacHash.Write([]byte(stringToSign))

	// The Base64 encoded HMAC hash is the signature of the request
	return base64.StdEncoding.EncodeToString(hmacHash.Sum(nil))
}

//escapePlusAndStarAndUnescapeTildeSymbol returns a string,
//in which the character "+" is escaped to "%20", character "*" is escaped to "%2A" and unescapes "%7E" to "~".
func escapePlusAndStarAndUnescapeTildeSymbol(str string) string {
	str = strings.ReplaceAll(str, "+", "%20")
	str = strings.ReplaceAll(str, "*", "%2A")
	str = strings.ReplaceAll(str, "%7E", "~")

	return str
}
