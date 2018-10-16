package model

type system struct {
	IPAddrAPI string `yaml:"IPAddrAPI"`
}

type cloudFlare struct {
	APIEndpoint string `yaml:"APIEndpoint"`
	APIKey string `yaml:"APIKey"`
	ZoneID string `yaml:"ZoneID"`
	AuthEmail string `yaml:"AuthEmail"`
	DomainName string `yaml:"DomainName"`
}

type Config struct {
	System system `yaml:"System"`
	CloudFlare cloudFlare `yaml:"CloudFlare"`
}
