package model

type DNSRecord struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

type SetDNSRecordRequest struct {
	IPs []string `json:"ips"`
}
