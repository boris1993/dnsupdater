package common

import "testing"

func TestCompareAddresses(t *testing.T) {
	var sameAddress bool
	var failTestMessage = "%s and %s should be the same address"

	sameAddress = CompareAddresses("192.168.1.1", "192.168.001.001")
	if !sameAddress {
		t.Errorf(failTestMessage, "192.168.1.1", "192.168.001.001")
	}

	sameAddress = CompareAddresses(
		"2001:0db8:0000:0000:0000:ff00:0042:8329", "2001:db8:0:0:0:ff00:42:8329")
	if !sameAddress {
		t.Errorf(failTestMessage, "2001:0db8:0000:0000:0000:ff00:0042:8329", "2001:db8:0:0:0:ff00:42:8329")
	}

	sameAddress = CompareAddresses(
		"2001:0db8:0000:0000:0000:ff00:0042:8329", "2001:db8::ff00:42:8329")
	if !sameAddress {
		t.Errorf(failTestMessage, "2001:0db8:0000:0000:0000:ff00:0042:8329", "2001:db8::ff00:42:8329")
	}

	sameAddress = CompareAddresses(
		"2001:db8:0:0:0:ff00:42:8329", "2001:db8::ff00:42:8329")
	if !sameAddress {
		t.Errorf(failTestMessage, "2001:db8:0:0:0:ff00:42:8329", "2001:db8::ff00:42:8329")
	}

	sameAddress = CompareAddresses("not valid address", "not valid address")
	if sameAddress {
		t.Error("Should return false when comparing invalid IP addresses")
	}
}
