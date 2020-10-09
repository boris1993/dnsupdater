package aliyun

type commonResponse struct {
	RequestID string `json:"RequestId"`
}

// DescribeDomainRecordsResponse corresponds to the return value of the DescribeDomainRecords API
// See document: https://help.aliyun.com/document_detail/29776.html?spm=a2c4g.11186623.2.37.1de5425696HU8m#h2-u8FD4u56DEu53C2u65703
type DescribeDomainRecordsResponse struct {
	commonResponse
	TotalCount    int32           `json:"TotalCount"`
	PageNumber    int32           `json:"PageNumber"`
	PageSize      int32           `json:"PageSize"`
	DomainRecords []domainRecords `json:"DomainRecords"`
}

// domainRecords describes how a DNS record is returned.
// See document: https://help.aliyun.com/document_detail/29799.html?spm=a2c4g.11186623.2.18.37a25eb4Fu6boQ
type domainRecords struct {
	DomainName string
	RecordID   string `json:"RecordId"`
	RR         string
	Type       string
	Value      string
	TTL        int32
	Priority   int32
	Line       string
	Status     string
	Locked     string
	Weight     int32
}
