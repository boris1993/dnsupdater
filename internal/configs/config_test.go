package configs

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
	Path = testResourcePath + "/test_config.yaml"
	if _, err := os.Stat(Path); os.IsNotExist(err) {
		t.Errorf("test_config.yaml doesn't exist")
		return
	}

	config, err := Get()
	if err != nil {
		t.Error(err)
		return
	}

	if config.System.IPAddrAPI == "" || config.System.CloudFlareAPIEndpoint == "" || config.System.AliyunAPIEndpoint == "" {
		t.Errorf("Content empty in the System part of test_config.yaml")
		return
	}

	if len(config.CloudFlareRecords) != 4 {
		t.Errorf("Error reading the CloudFlareRecords part. Expected 4 records but found %d", len(config.CloudFlareRecords))
		return
	}

	if len(config.AliDNSRecords) != 3 {
		t.Errorf("Error reading the AliDNSRecords part. Expected 2 records but found %d", len(config.AliDNSRecords))
		return
	}
}

func testGetFail(t *testing.T) {
	Debug = true
	Path = testResourcePath + "/non_existent_config.yaml"

	_, err := Get()
	if err == nil {
		t.Error("TestGetFail should fail")
	}
}
