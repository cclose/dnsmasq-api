package service

import (
	"github.com/cclose/dnsmasq-api/model"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"reflect"
	"testing"
)

func TestDNSMasqService_BuildDatabase(t *testing.T) {
	type fields struct {
		db                *bolt.DB
		dnsBucket         []byte
		dbFilePath        string
		dnsMasqConfig     string
		skipDNSMasqReload bool
		log               *logrus.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DNSMasqService{
				db:                tt.fields.db,
				dnsBucket:         tt.fields.dnsBucket,
				dbFilePath:        tt.fields.dbFilePath,
				dnsMasqConfig:     tt.fields.dnsMasqConfig,
				skipDNSMasqReload: tt.fields.skipDNSMasqReload,
				log:               tt.fields.log,
			}
			if err := ds.BuildDatabase(); (err != nil) != tt.wantErr {
				t.Errorf("BuildDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDNSMasqService_DeleteByHost(t *testing.T) {
	type fields struct {
		db                *bolt.DB
		dnsBucket         []byte
		dbFilePath        string
		dnsMasqConfig     string
		skipDNSMasqReload bool
		log               *logrus.Logger
	}
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DNSMasqService{
				db:                tt.fields.db,
				dnsBucket:         tt.fields.dnsBucket,
				dbFilePath:        tt.fields.dbFilePath,
				dnsMasqConfig:     tt.fields.dnsMasqConfig,
				skipDNSMasqReload: tt.fields.skipDNSMasqReload,
				log:               tt.fields.log,
			}
			if err := ds.DeleteByHost(tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("DeleteByHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDNSMasqService_GetAllIPs(t *testing.T) {
	type fields struct {
		db                *bolt.DB
		dnsBucket         []byte
		dbFilePath        string
		dnsMasqConfig     string
		skipDNSMasqReload bool
		log               *logrus.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		want    []model.DNSRecord
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DNSMasqService{
				db:                tt.fields.db,
				dnsBucket:         tt.fields.dnsBucket,
				dbFilePath:        tt.fields.dbFilePath,
				dnsMasqConfig:     tt.fields.dnsMasqConfig,
				skipDNSMasqReload: tt.fields.skipDNSMasqReload,
				log:               tt.fields.log,
			}
			got, err := ds.GetAllIPs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllIPs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllIPs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDNSMasqService_GetIPByHost(t *testing.T) {
	type fields struct {
		db                *bolt.DB
		dnsBucket         []byte
		dbFilePath        string
		dnsMasqConfig     string
		skipDNSMasqReload bool
		log               *logrus.Logger
	}
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.DNSRecord
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DNSMasqService{
				db:                tt.fields.db,
				dnsBucket:         tt.fields.dnsBucket,
				dbFilePath:        tt.fields.dbFilePath,
				dnsMasqConfig:     tt.fields.dnsMasqConfig,
				skipDNSMasqReload: tt.fields.skipDNSMasqReload,
				log:               tt.fields.log,
			}
			got, err := ds.GetIPByHost(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIPByHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIPByHost() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDNSMasqService_ReloadDNSMasq(t *testing.T) {
	type fields struct {
		db                *bolt.DB
		dnsBucket         []byte
		dbFilePath        string
		dnsMasqConfig     string
		skipDNSMasqReload bool
		log               *logrus.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DNSMasqService{
				db:                tt.fields.db,
				dnsBucket:         tt.fields.dnsBucket,
				dbFilePath:        tt.fields.dbFilePath,
				dnsMasqConfig:     tt.fields.dnsMasqConfig,
				skipDNSMasqReload: tt.fields.skipDNSMasqReload,
				log:               tt.fields.log,
			}
			if err := ds.ReloadDNSMasq(); (err != nil) != tt.wantErr {
				t.Errorf("ReloadDNSMasq() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDNSMasqService_SetIPByHost(t *testing.T) {
	type fields struct {
		db                *bolt.DB
		dnsBucket         []byte
		dbFilePath        string
		dnsMasqConfig     string
		skipDNSMasqReload bool
		log               *logrus.Logger
	}
	type args struct {
		hostname string
		ips      []string
		appendIP bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.DNSRecord
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DNSMasqService{
				db:                tt.fields.db,
				dnsBucket:         tt.fields.dnsBucket,
				dbFilePath:        tt.fields.dbFilePath,
				dnsMasqConfig:     tt.fields.dnsMasqConfig,
				skipDNSMasqReload: tt.fields.skipDNSMasqReload,
				log:               tt.fields.log,
			}
			got, err := ds.SetIPByHost(tt.args.hostname, tt.args.ips, tt.args.appendIP)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetIPByHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetIPByHost() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDNSMasqService_UpdateDNSMasq(t *testing.T) {
	type fields struct {
		db                *bolt.DB
		dnsBucket         []byte
		dbFilePath        string
		dnsMasqConfig     string
		skipDNSMasqReload bool
		log               *logrus.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DNSMasqService{
				db:                tt.fields.db,
				dnsBucket:         tt.fields.dnsBucket,
				dbFilePath:        tt.fields.dbFilePath,
				dnsMasqConfig:     tt.fields.dnsMasqConfig,
				skipDNSMasqReload: tt.fields.skipDNSMasqReload,
				log:               tt.fields.log,
			}
			if err := ds.UpdateDNSMasq(); (err != nil) != tt.wantErr {
				t.Errorf("UpdateDNSMasq() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDNSMasqService_WriteDNSMasq(t *testing.T) {
	type fields struct {
		db                *bolt.DB
		dnsBucket         []byte
		dbFilePath        string
		dnsMasqConfig     string
		skipDNSMasqReload bool
		log               *logrus.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DNSMasqService{
				db:                tt.fields.db,
				dnsBucket:         tt.fields.dnsBucket,
				dbFilePath:        tt.fields.dbFilePath,
				dnsMasqConfig:     tt.fields.dnsMasqConfig,
				skipDNSMasqReload: tt.fields.skipDNSMasqReload,
				log:               tt.fields.log,
			}
			if err := ds.WriteDNSMasq(); (err != nil) != tt.wantErr {
				t.Errorf("WriteDNSMasq() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDNSMasqService_openDB(t *testing.T) {
	type fields struct {
		db                *bolt.DB
		dnsBucket         []byte
		dbFilePath        string
		dnsMasqConfig     string
		skipDNSMasqReload bool
		log               *logrus.Logger
	}
	type args struct {
		dbPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DNSMasqService{
				db:                tt.fields.db,
				dnsBucket:         tt.fields.dnsBucket,
				dbFilePath:        tt.fields.dbFilePath,
				dnsMasqConfig:     tt.fields.dnsMasqConfig,
				skipDNSMasqReload: tt.fields.skipDNSMasqReload,
				log:               tt.fields.log,
			}
			if err := ds.openDB(tt.args.dbPath); (err != nil) != tt.wantErr {
				t.Errorf("openDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDNSMasqService(t *testing.T) {
	type args struct {
		config model.Config
		opts   []DNSMasqServiceOption
	}
	tests := []struct {
		name    string
		args    args
		want    IDNSMasqService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDNSMasqService(tt.args.config, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDNSMasqService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDNSMasqService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDBFilePath(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want DNSMasqServiceOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDBFilePath(tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDBFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDNSBucket(t *testing.T) {
	type args struct {
		bucket string
	}
	tests := []struct {
		name string
		args args
		want DNSMasqServiceOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDNSBucket(tt.args.bucket); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDNSBucket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	type args struct {
		logger *logrus.Logger
	}
	tests := []struct {
		name string
		args args
		want DNSMasqServiceOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithLogger(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeDuplicates(t *testing.T) {
	type args struct {
		records []model.DNSRecord
	}
	tests := []struct {
		name string
		args args
		want []model.DNSRecord
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDuplicates(tt.args.records); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}
