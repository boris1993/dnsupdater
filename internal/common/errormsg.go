// Package common contains all constants needed in this programme
package common

const IPv4 = "IPv4"
const IPv6 = "IPv6"

const MsgHeaderDNSRecordUpdateSuccessful = "Successfully updated the DNS record"
const MsgHeaderCurrentIPv4Addr = "Current IPv4 address is:"
const MsgHeaderCurrentIPv6Addr = "Current IPv6 address is:"
const MsgIPv4Disabled = "IPv4 disabled. Will skip checking and updating all the A records"
const MsgIPv6Disabled = "IPv6 disabled. Will skip checking and updating all the AAAA records"
const MsgIPv4AddrNotAvailable = "No valid IPv4 internet address. Skipping updating this record."
const MsgIPv6AddrNotAvailable = "No valid IPv6 internet address. Skipping updating this record."
const MsgHeaderFetchingIPOfDomain = "Fetching the IP address of domain"
const MsgHeaderLoadingConfig = "Loading configurations from"
const MsgTemplateDomainProcessing = "Processing %s, type %s"
const MsgCloudFlareRecordsFoundSuffix = "CloudFlare DNS record(s) found"
const MsgAliDNSRecordsFoundSuffix = "Aliyun DNS record(s) found"

const MsgFormatDNSFetchResult = "The IP address of %s is %s"
const MsgFormatAliDNSFetchResult = "The IP address of %s is %s, and the record ID is %s"
const MsgFormatUpdatingDNS = "Updating the IP address of record ID %s to %s"

const ErrMsgHeaderFetchDomainInfoFailed = "Failed to get the information for domain"
const ErrMsgHeaderUpdateDNSRecordFailed = "Failed to update the DNS record"
const MsgIPAddrNotChanged = "IP address not changed. Will not update the DNS record."

const MsgCheckingCurrentIPv4Addr = "Checking current IPv4 address..."
const MsgCheckingCurrentIPv6Addr = "Checking current IPv6 address..."
const ErrNoDNSRecordFoundPrefix = "No record matches the domain name "
const ErrInvalidDomainType = "The type of the domain is invalid. Only A and AAAA are accepted."

const ErrCloseHTTPConnectionFail = "Error closing the HTTP connection"
const ErrIPAddressFetchingAPIEmpty = "API address for fetching current IP address cannot be empty"
const ErrCloudFlareAPIAddressEmpty = "CloudFlare API endpoint address cannot be empty"
const ErrAliyunAPIAddressEmpty = "Aliyun API endpoint address cannot be empty"
const ErrCloudFlareRecordConfigIncomplete = "Incomplete CloudFlare configuration found"
const ErrAliDNSRecordConfigIncomplete = "Incomplete Aliyun DNS configuration found"
const ErrJsonPathNotSpecified = "JSON Path for %s is not specified"
const ErrInvalidResponseTypeSpecified = "The specified response type for %s is invalid"

const WarnAuthEmailDeprecated = "AuthEmail is deprecated. Please use dedicated API token if you are still using Global API Key."
