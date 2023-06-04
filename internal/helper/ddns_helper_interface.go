package helper

type DDNSHelperInterface interface {
	ProcessRecords(currentIPv4Address string, currentIPv6Address string) error
}
