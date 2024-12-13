package model

import "time"

type AppConfig struct {
	Config    Config
	BuildInfo BuildInfo
}

type Config struct {
	DnsmasqConfig     string         `mapstructure:"dnsmasq_config"`
	DB                DatabaseConfig `mapstructure:"db"`
	Logging           LoggingConfig  `mapstructure:"logging"`
	Port              int            `mapstructure:"port"`
	SkipDNSMasqReload bool           `mapstructure:"skip_dnsmasq_reload"`
	SSL               SSLConfig      `mapstructure:"ssl"`
}

type DatabaseConfig struct {
	FilePath   string `mapstructure:"file_path"`
	BucketName string `mapstructure:"bucket_name"`
}

type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

type SSLConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

type BuildInfo struct {
	BuildTime time.Time
	Commit    string
	Version   string
}
