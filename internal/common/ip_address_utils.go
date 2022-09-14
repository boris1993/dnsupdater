package common

import (
	"net"
	"strings"
)

// CompareAddresses compares 2 given IP address string and see if they are the same IP address
func CompareAddresses(address1 string, address2 string) bool {
	ipAddr1 := net.ParseIP(removeLeadingZerosFromIPv4Address(address1))
	ipAddr2 := net.ParseIP(removeLeadingZerosFromIPv4Address(address2))

	if ipAddr1 == nil || ipAddr2 == nil {
		return false
	}

	return ipAddr1.Equal(ipAddr2)
}

func removeLeadingZerosFromIPv4Address(ipAddress string) string {
	ipAddressParts := strings.Split(ipAddress, ".")

	for index, part := range ipAddressParts {
		ipAddressParts[index] = strings.TrimLeft(part, "0")
	}

	return strings.Join(ipAddressParts, ".")
}
