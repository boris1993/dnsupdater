// Package constants contains all constants needed in this programme
package constants

const MsgHeaderDNSRecordUpdateSuccessful = "Successfully updated the DNS record"
const MsgHeaderCurrentIPAddr = "Current IP address is:"
const MsgHeaderFetchingIPOfDomain = "Fetching the IP address of domain"
const MsgHeaderLoadingConfig = "Loading configuration from"
const MsgHeaderDomainProcessing = "Processing"
const MsgCloudFlareRecordsFoundSuffix = "CloudFlare DNS record(s) found"
const MsgAliDNSRecordsFoundSuffix = "Aliyun DNS record(s) found"

const MsgFormatDNSFetchResult = "The IP address of %s is %s"
const MsgFormatAliDNSFetchResult = "The IP address of %s is %s, and the record ID is %s"
const MsgFormatUpdatingDNS = "Updating the IP address of record ID %s to %s"

const ErrMsgHeaderFetchDomainInfoFailed = "Failed to get the information for domain"
const ErrMsgHeaderUpdateDNSRecordFailed = "Failed to update the DNS record"
const MsgIPAddrNotChanged = "IP address not changed. Will not update the DNS record."

const MsgCheckingCurrentIPAddr = "Checking current IP address..."
const ErrNoDNSRecordFoundPrefix = "No record matches the domain name "
const ErrMultipleDNSRecordsFound = "Found multiple records"

const ErrCloseHTTPConnectionFail = "Error closing the HTTP connection"
const ErrIPAddressFetchingAPIEmpty = "API address for fetching current IP address cannot be empty"
const ErrCloudFlareAPIAddressEmpty = "CloudFlare API endpoint address cannot be empty"
const ErrCloudFlareRecordConfigIncomplete = "Incomplete CloudFlare configuration found"
const ErrAliDNSRecordConfigIncomplete = "Incomplete Aliyun DNS configuration found"
