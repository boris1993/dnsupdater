// Package model provides all models needed by this programme.
package model

// CfDnsRecord is the structure of the CloudFlare DNS query.
//
// See https://api.cloudflare.com/#dns-records-for-a-zone-dns-record-details
type CfDnsRecord struct {
	Success    bool                `json:"success"`
	Errors     []string            `json:"errors"`
	Messages   []string            `json:"messages"`
	Result     []CfDnsRecordResult `json:"result"`
	ResultInfo []interface{}       `json:"result_info"`
}

// CfDnsRecordResult is the structure of the results part in CfDnsRecord.
//
// See https://api.cloudflare.com/#dns-records-for-a-zone-dns-record-details
type CfDnsRecordResult struct {
	ID         string        `json:"id"`
	RecordType string        `json:"type"`
	Name       string        `json:"name"`
	Content    string        `json:"content"`
	Proxiable  string        `json:"proxiable"`
	Proxied    string        `json:"proxied"`
	TTL        int           `json:"ttl"`
	Locked     bool          `json:"locked"`
	ZoneID     string        `json:"zone_id"`
	ZoneName   string        `json:"zone_name"`
	CreatedOn  string        `json:"created_on"`
	ModifiedOn string        `json:"modified_on"`
	Data       []interface{} `json:"data"`
}

// UpdateRecordData describes the required parameters when updating a DNS record.
//
// See https://api.cloudflare.com/#dns-records-for-a-zone-update-dns-record
type UpdateRecordData struct {
	RecordType string `json:"type"`
	Name       string `json:"name"`
	Content    string `json:"content"`
}

// UpdateRecordResult describes what will be returned
// when updating a DNS record.
//
// See https://api.cloudflare.com/#dns-records-for-a-zone-update-dns-record
type UpdateRecordResult struct {
	Success  bool              `json:"success"`
	Errors   []string          `json:"errors"`
	Messages []string          `json:"messages"`
	Result   CfDnsRecordResult `json:"result"`
}
