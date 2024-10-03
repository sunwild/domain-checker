package domains

// strucrture for being "Domain"
type Domain struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DomainStatus struct {
	Domain     string `json:"domain"`
	StatusCode int    `json:"http_status"`
	DnsStatus  string `json:"dns_status"`
	SslStatus  string `json:"ssl_status"`
	VirusTotal string `json:"virus_total_status"`
}