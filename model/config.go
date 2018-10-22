// Package model provides all models needed by this programme.
package model

// system describes the system properties in config.yaml
type system struct {
	IPAddrAPI string `yaml:"IPAddrAPI"`
}

// cloudFlare describes the CloudFlare specialised properties in config.yaml
type cloudFlare struct {
	APIEndpoint string `yaml:"APIEndpoint"`
	APIKey      string `yaml:"APIKey"`
	ZoneID      string `yaml:"ZoneID"`
	AuthEmail   string `yaml:"AuthEmail"`
	DomainName  string `yaml:"DomainName"`
}

// Config describes the top-level properties in config.yaml
type Config struct {
	System     system     `yaml:"System"`
	CloudFlare cloudFlare `yaml:"CloudFlare"`
}
