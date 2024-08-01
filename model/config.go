package model

import "time"

type AppConfig struct {
	Config    Config
	BuildInfo BuildInfo
}

type Config struct {
	Port              int    `mapstructure:"port"`
	SSLEnabled        bool   `mapstructure:"ssl.enabled"`
	SSLCertFile       string `mapstructure:"ssl.cert_file"`
	SSLKeyFile        string `mapstructure:"ssl.key_file"`
	DnsmasqConfig     string `mapstructure:"dnsmasq_config"`
	SkipDNSMasqReload bool   `mapstructure:"skip_dnsmasq_reload"`
}

type BuildInfo struct {
	BuildTime time.Time
	Commit    string
	Version   string
}
