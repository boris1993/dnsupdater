package conf

import (
	"os"
	"sync"
	"testing"
)

const testResourcePath = "../test"

func TestGet(t *testing.T) {
	t.Run("TestGetSuccess", testGetSuccess)
	once = new(sync.Once)
	t.Run("testGetFail", testGetFail)
}

func testGetSuccess(t *testing.T) {
	Debug = true
	ConfigFilePath = testResourcePath + "/test_config.yaml"
	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		t.Errorf("test_config.yaml doesn't exist")
		return
	}

	config, err := GetConfig()
	if err != nil {
		t.Error(err)
		return
	}

	if config.System.IPv4.IPAddrAPI == "" ||
		config.System.IPv6.IPAddrAPI == "" ||
		config.System.Endpoints.CloudFlareAPIEndpoint == "" ||
		config.System.Endpoints.AliyunAPIEndpoint == "" {
		t.Errorf("Content empty in the System part of test_config.yaml")
		return
	}

	if len(config.CloudFlareRecords) != 6 {
		t.Errorf("Error reading the CloudFlareRecords part. Expected 4 records but found %d", len(config.CloudFlareRecords))
		return
	}

	if len(config.AliDNSRecords) != 5 {
		t.Errorf("Error reading the AliDNSRecords part. Expected 5 records but found %d", len(config.AliDNSRecords))
		return
	}
}

func testGetFail(t *testing.T) {
	Debug = true
	ConfigFilePath = testResourcePath + "/non_existent_config.yaml"

	_, err := GetConfig()
	if err == nil {
		t.Error("TestGetFail should fail")
	}
}
