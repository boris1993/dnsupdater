package model

type CfDnsRecord struct {
	Success    bool                `json:"success"`
	Errors     []string            `json:"errors"`
	Messages   []string            `json:"messages"`
	Result     []CfDnsRecordResult `json:"result"`
	ResultInfo []interface{}       `json:"result_info"`
}

type CfDnsRecordResult struct {
	Id         string        `json:"id"`
	RecordType string        `json:"type"`
	Name       string        `json:"name"`
	Content    string        `json:"content"`
	Proxiable  string        `json:"proxiable"`
	Proxied    string        `json:"proxied"`
	Ttl        int           `json:"ttl"`
	Locked     bool          `json:"locked"`
	ZoneId     string        `json:"zone_id"`
	ZoneName   string        `json:"zone_name"`
	CreatedOn  string        `json:"created_on"`
	ModifiedOn string        `json:"modified_on"`
	Data       []interface{} `json:"data"`
}

type UpdateRecordData struct {
	RecordType string `json:"type"`
	Name       string `json:"name"`
	Content    string `json:"content"`
}

type UpdateRecordResult struct {
	Success  bool              `json:"success"`
	Errors   []string          `json:"errors"`
	Messages []string          `json:"messages"`
	Result   CfDnsRecordResult `json:"result"`
}
