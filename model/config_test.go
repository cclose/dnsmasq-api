package model

import (
	"github.com/mitchellh/mapstructure"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfig_ViperLoad(t *testing.T) {
	tests := []struct {
		name              string
		args              string
		want              Config
		wantParseErr      bool
		wantUnmarshallErr error
	}{
		{
			name: "ValidConfig",
			args: `---
port: 8080
ssl:
  enabled: true
  cert_file: "/path/to/cert"
  key_file: "/path/to/key"
dnsmasq_config: "/path/to/dnsmasq.conf"
skip_dnsmasq_reload: true
db:
  file_path: "/path/to/db"
  bucket_name: "mybucket"
logging:
  level: "info"
  output: "stdout"
  file_path: ""
`,
			want: Config{
				DnsmasqConfig: "/path/to/dnsmasq.conf",
				DB: DatabaseConfig{
					FilePath:   "/path/to/db",
					BucketName: "mybucket",
				},
				Logging: LoggingConfig{
					Level:    "info",
					Output:   "stdout",
					FilePath: "",
				},
				Port:              8080,
				SkipDNSMasqReload: true,
				SSL: SSLConfig{
					Enabled:  true,
					CertFile: "/path/to/cert",
					KeyFile:  "/path/to/key",
				},
			},
		},
		{
			name: "MissingRequiredFields",
			want: Config{
				DnsmasqConfig: "",
				DB: DatabaseConfig{
					FilePath:   "",
					BucketName: "",
				},
				Logging: LoggingConfig{
					Level:    "",
					Output:   "",
					FilePath: "",
				},
				Port:              0,
				SkipDNSMasqReload: false,
				SSL: SSLConfig{
					Enabled:  false,
					CertFile: "",
					KeyFile:  "",
				},
			},
		},
		{
			name: "config.yaml.dist",
			args: `---
port: 8080
ssl:
  enabled: false
dnsmasq_config: "/etc/dnsmasq.d/api.conf"
skip_dnsmasq_reload: true`,
			want: Config{
				Port:              8080,
				SkipDNSMasqReload: true,
				DnsmasqConfig:     "/etc/dnsmasq.d/api.conf",
				SSL:               SSLConfig{Enabled: false},
			},
		},
		{
			name: "Invalid Boolean Value",
			args: `---
skip_dnsmasq_reload: flase`,
			wantUnmarshallErr: &mapstructure.Error{Errors: []string{`cannot parse 'skip_dnsmasq_reload' as bool: strconv.ParseBool: parsing "flase": invalid syntax`}},
		},
		{
			name: "Invalid int value",
			args: `---
port: fish
`,
			wantUnmarshallErr: &mapstructure.Error{Errors: []string{`cannot parse 'port' as int: strconv.ParseInt: parsing "fish": invalid syntax`}},
		},
		{
			name:         "Invalid Config",
			args:         "\n\t\t\t\tFish = DUmplings {{}}",
			wantParseErr: true,
		},
		{
			name: "JSON",
			args: `{"port": 8080, "ssl":{"enabled": false, "cert_file": "/path/to/cert"}}`,
			want: Config{
				Port: 8080,
				SSL: SSLConfig{
					Enabled:  false,
					CertFile: "/path/to/cert",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			v.SetConfigType("yaml")
			err := v.ReadConfig(strings.NewReader(tt.args))
			assert.Equal(t, tt.wantParseErr, err != nil)

			var config Config
			err = v.Unmarshal(&config)
			assert.Equal(t, tt.wantUnmarshallErr, err)

			assert.Equal(t, tt.want.Port, config.Port)
			assert.Equal(t, tt.want.SSL.Enabled, config.SSL.Enabled)
			assert.Equal(t, tt.want.SSL.CertFile, config.SSL.CertFile)
			assert.Equal(t, tt.want.SSL.KeyFile, config.SSL.KeyFile)
			assert.Equal(t, tt.want.DnsmasqConfig, config.DnsmasqConfig)
			assert.Equal(t, tt.want.SkipDNSMasqReload, config.SkipDNSMasqReload)
			assert.Equal(t, tt.want.DB.FilePath, config.DB.FilePath)
			assert.Equal(t, tt.want.DB.BucketName, config.DB.BucketName)
			assert.Equal(t, tt.want.Logging.Level, config.Logging.Level)
			assert.Equal(t, tt.want.Logging.Output, config.Logging.Output)
			assert.Equal(t, tt.want.Logging.FilePath, config.Logging.FilePath)
		})
	}
}
