package common

import "net"

// CompareAddresses compares 2 given IP address string and see if they are the same IP address
func CompareAddresses(address1 string, address2 string) bool {
	ipAddr1 := net.ParseIP(address1)
	ipAddr2 := net.ParseIP(address2)

	if ipAddr1 == nil || ipAddr2 == nil {
		return false
	}

	return ipAddr1.Equal(ipAddr2)
}
