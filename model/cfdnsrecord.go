// Package model provides all models needed by this programme.
package model

// CfDnsRecord is the structure of the CloudFlare DNS query.
//
// See https://api.cloudflare.com/#dns-records-for-a-zone-dns-record-details
type CfDnsRecord struct {
	Success    bool                  `json:"success"`
	Errors     []string              `json:"errors"`
	Messages   []string              `json:"messages"`
	Result     []cfDnsRecordResult   `json:"result"`
	ResultInfo cfDnsRecordResultInfo `json:"result_info"`
}

// cfDnsRecordResult is the structure of the results part in CfDnsRecord.
//
// See https://api.cloudflare.com/#dns-records-for-a-zone-dns-record-details
type cfDnsRecordResult struct {
	ID         string                `json:"id"`
	RecordType string                `json:"type"`
	Name       string                `json:"name"`
	Content    string                `json:"content"`
	Proxiable  bool                  `json:"proxiable"`
	Proxied    bool                  `json:"proxied"`
	TTL        int                   `json:"ttl"`
	Locked     bool                  `json:"locked"`
	Meta       cfDnsRecordResultMeta `json:"meta"`
	ZoneID     string                `json:"zone_id"`
	ZoneName   string                `json:"zone_name"`
	CreatedOn  string                `json:"created_on"`
	ModifiedOn string                `json:"modified_on"`
	Data       []interface{}         `json:"data"`
}

type cfDnsRecordResultInfo struct {
	Count      int `json:"count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

type cfDnsRecordResultMeta struct {
	AutoAdded           bool `json:"auto_added"`
	ManagedByApps       bool `json:"managed_by_apps"`
	ManagedByArgoTunnel bool `json:"managed_by_argo_tunnel"`
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
	Result   cfDnsRecordResult `json:"result"`
}
