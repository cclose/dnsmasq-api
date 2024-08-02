package service

import (
	"encoding/json"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/cclose/dnsmasq-api/model"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"os"
	"os/exec"
	"strings"
)

const (
	dbFileMode          os.FileMode = 0600
	dnsFileMode         os.FileMode = 0644
	dnsConfigHeader                 = "# Managed by DNSMasq API\n"
	defaultDBFilePath               = "dns.db"
	defaultDBBucketName             = "dnsRecords"

	MetricDNSCount   = "dnsmasq_hostname_total"
	MetricIPCount    = "dnsmasq_ip_total"
	MetricDNSReloads = "dnsmasq_reloads_total"
)

type IDNSMasqService interface {
	BuildDatabase() error
	ReloadDNSMasq() error
	UpdateDNSMasq() error
	WriteDNSMasq() error

	GetAllIPs() ([]model.DNSRecord, error)
	GetIPByHost(host string) ([]model.DNSRecord, error)
	SetIPByHost(hostname string, ips []string, appendIP bool) ([]model.DNSRecord, error)
	DeleteByHost(host string) error
}

type DNSMasqService struct {
	db *bolt.DB

	dnsBucket  []byte
	dbFilePath string

	dnsMasqConfig     string
	skipDNSMasqReload bool

	log *logrus.Logger
}

// DNSMasqServiceOption Option functions for customizing DNSMasqService from Constructor
type DNSMasqServiceOption func(*DNSMasqService)

// NewDNSMasqService Creates a new DNSMasqService
func NewDNSMasqService(config model.Config, opts ...DNSMasqServiceOption) (IDNSMasqService, error) {
	ds := &DNSMasqService{
		dbFilePath: defaultDBFilePath,
		dnsBucket:  []byte(defaultDBBucketName),

		dnsMasqConfig:     config.DnsmasqConfig,
		skipDNSMasqReload: config.SkipDNSMasqReload,
	}

	// Apply any options
	for _, opt := range opts {
		opt(ds)
	}

	if ds.log == nil {
		ds.log = logrus.New()
	}

	if err := ds.openDB(ds.dbFilePath); err != nil {
		return nil, err
	}

	if err := ds.BuildDatabase(); err != nil {
		return nil, err
	}

	return ds, nil
}

// Option Functions

// WithDNSBucket Sets the name of the DB Bucket to store DNS Records
func WithDNSBucket(bucket string) DNSMasqServiceOption {
	return func(ds *DNSMasqService) {
		ds.dnsBucket = []byte(bucket)
	}
}

// WithLogger Sets the logger for the service to use
func WithLogger(logger *logrus.Logger) DNSMasqServiceOption {
	return func(ds *DNSMasqService) {
		ds.log = logger
	}
}

// WithDBFilePath Sets the path to the db file for storing DNS Records
func WithDBFilePath(filePath string) DNSMasqServiceOption {
	return func(ds *DNSMasqService) {
		ds.dbFilePath = filePath
	}
}

// WithConfig Creates DNSMasqServiceOptions from a DatabaseConfig
func WithConfig(dbConfig model.DatabaseConfig) DNSMasqServiceOption {
	options := []DNSMasqServiceOption{}
	if dbConfig.BucketName != "" && dbConfig.BucketName != defaultDBBucketName {
		options = append(options, WithDNSBucket(dbConfig.BucketName))
	}
	if dbConfig.FilePath != "" && dbConfig.FilePath != defaultDBFilePath {
		options = append(options, WithDBFilePath(dbConfig.FilePath))
	}

	return func(ds *DNSMasqService) {
		for _, option := range options {
			option(ds)
		}
	}
}

func (ds *DNSMasqService) openDB(dbPath string) (err error) {
	ds.db, err = bolt.Open(dbPath, dbFileMode, nil)
	if err != nil {
		return
	}

	// Make sure our DNS Bucket exists
	err = ds.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(ds.dnsBucket)
		return err
	})

	return
}

// GetAllIPs retrieves all DNS records from the database.
func (ds *DNSMasqService) GetAllIPs() ([]model.DNSRecord, error) {
	var records []model.DNSRecord

	err := ds.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(ds.dnsBucket)
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		return bucket.ForEach(func(k, v []byte) error {
			var hostRecords []model.DNSRecord
			if err := json.Unmarshal(v, &hostRecords); err != nil {
				return err
			}
			// append all host records to records
			records = append(records, hostRecords...)
			return nil
		})
	})

	return records, err
}

const ErrorNoIPForHost = "no records found for host"

// GetIPByHost retrieves all IP addresses for the given hostname.
func (ds *DNSMasqService) GetIPByHost(host string) ([]model.DNSRecord, error) {
	var records []model.DNSRecord

	err := ds.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(ds.dnsBucket)
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		data := bucket.Get([]byte(host))
		if data == nil {
			return fmt.Errorf(ErrorNoIPForHost) // Return error if no records found
		}

		if err := json.Unmarshal(data, &records); err != nil {
			return err
		}

		return nil
	})

	return records, err
}

// removeDuplicates removes duplicate DNS records based on IP.
func removeDuplicates(records []model.DNSRecord) []model.DNSRecord {
	seen := make(map[string]bool)
	var uniqueRecords []model.DNSRecord

	for _, record := range records {
		if !seen[record.IP] {
			seen[record.IP] = true
			uniqueRecords = append(uniqueRecords, record)
		}
	}

	return uniqueRecords
}

// SetIPByHost sets or appends an IP address for the given hostname.
// If appendIP is true, it will add the IP to the existing list, otherwise it will replace it.
func (ds *DNSMasqService) SetIPByHost(hostname string, ips []string, appendIP bool) ([]model.DNSRecord, error) {
	var records []model.DNSRecord
	err := ds.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(ds.dnsBucket)
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		for _, ip := range ips {
			records = append(records, model.DNSRecord{
				Hostname: hostname, IP: ip,
			})
		}

		if appendIP {
			currRecords, err := ds.GetIPByHost(hostname)
			if err != nil && err.Error() != ErrorNoIPForHost {
				return err
			}
			records = append(currRecords, records...)
		}

		records = removeDuplicates(records)

		newData, err := json.Marshal(records)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(hostname), newData)
	})

	return records, err
}

// DeleteByHost deletes all IP addresses for the given hostname.
func (ds *DNSMasqService) DeleteByHost(host string) error {
	return ds.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(ds.dnsBucket)
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		return bucket.Delete([]byte(host))
	})
}

// BuildDatabase reads the DNSMasq config file and syncs the in-memory database.
func (ds *DNSMasqService) BuildDatabase() error {
	dnsCount := 0
	ipsCount := 0

	data, err := os.ReadFile(ds.dnsMasqConfig)
	if err != nil {
		ds.log.Fatalf("Failed to read dnsmasq config file: %v", err)
	}
	lines := strings.Split(string(data), "\n")
	entries := make(map[string][]string)
	for _, line := range lines {
		if strings.HasPrefix(line, "address=/") {
			parts := strings.Split(line, "/")
			if len(parts) == 3 {
				var ips []string
				var ok bool
				hostname := parts[1]
				ip := strings.TrimSpace(parts[2])
				if ips, ok = entries[hostname]; ok {
					ips = append(ips, ip)
				} else {
					ips = []string{ip}
				}
				entries[hostname] = ips
			}
		}
	}

	err = ds.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ds.dnsBucket)
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		err = b.ForEach(func(k, v []byte) error {
			return b.Delete(k)
		})

		return err
	})
	if err != nil {
		return err
	}

	for k, v := range entries {
		dnsCount += 1
		ipsCount += len(v)
		_, err = ds.SetIPByHost(k, v, false)
		if err != nil {
			return err
		}
	}

	metrics.GetOrCreateCounter(MetricDNSCount).Set(uint64(dnsCount))
	metrics.GetOrCreateCounter(MetricIPCount).Set(uint64(ipsCount))

	return nil
}

// ReloadDNSMasq Calls DNSMasq to reload it's config
func (ds *DNSMasqService) ReloadDNSMasq() error {
	metrics.GetOrCreateCounter(MetricDNSReloads).Inc()
	cmd := exec.Command("sudo", "systemctl", "reload", "dnsmasq")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// UpdateDNSMasq Syncs the in-memory DB to the DNS Masq file and reloads the service
func (ds *DNSMasqService) UpdateDNSMasq() error {
	err := ds.WriteDNSMasq()
	if err != nil {
		return err
	}

	if !ds.skipDNSMasqReload {
		err = ds.ReloadDNSMasq()
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteDNSMasq Writes the in-memory database out to the DNS Masq config file
func (ds *DNSMasqService) WriteDNSMasq() error {

	// Set our file header
	dnsConfigData := dnsConfigHeader
	ips, err := ds.GetAllIPs()
	if err != nil {
		return err
	}

	uniqHosts := make(map[string]struct{})

	for _, ip := range ips {
		// Track uniq hostnames
		uniqHosts[ip.Hostname] = struct{}{}
		dnsConfigData += fmt.Sprintf("address=/%s/%s\n", ip.Hostname, ip.IP)
	}

	// Write out the file
	err = os.WriteFile(ds.dnsMasqConfig, []byte(dnsConfigData), dnsFileMode)

	metrics.GetOrCreateCounter(MetricDNSCount).Set(uint64(len(uniqHosts)))
	metrics.GetOrCreateCounter(MetricIPCount).Set(uint64(len(ips)))

	return err
}
