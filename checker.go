package tldcheck

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/miekg/dns"
)

var (
	//nolint: gochecknoglobals
	DNSServers = []string{
		"1.1.1.1",
		"8.8.4.4",
		"8.8.8.8",
		"9.9.9.9",
	}
	//nolint: gochecknoglobals
	DNSPort = "53"

	//nolint: gochecknoglobals
	DNSTimeout = 8 * time.Second
	//nolint: gochecknoglobals
	DNSAttempts = 8
)

type checker struct {
	client *dns.Client
	config *dns.ClientConfig
	conn   *dns.Conn
	msg    *dns.Msg
}

func (c *checker) dnsQuery(domain string) (*dns.Msg, error) {
	c.msg.SetQuestion(dns.Fqdn(domain), dns.TypeNS)

	r, _, err := c.client.ExchangeWithConn(c.msg, c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", domain, err)
	}
	return r, nil //nolint: nlreturn
}

func (c *checker) Check(wi workItem) (*Result, error) {
	r, err := c.dnsQuery(wi.domain.DomainRaw())
	if err != nil {
		return nil, err
	}

	var available bool
	if r.Rcode == dns.RcodeNameError {
		available = true
	}

	return &Result{
		Domain:    wi.domain,
		Available: available,
	}, nil
}

func newChecker(c *dns.ClientConfig) (*checker, error) {
	//nolint: exhaustivestruct
	client := &dns.Client{
		Timeout: DNSTimeout,
	}
	server := c.Servers[rand.Intn(len(c.Servers))] //nolint: gosec
	port := c.Port
	socket := fmt.Sprintf("%s:%s", server, port)
	conn, err := client.Dial(socket)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", socket, err)
	}
	//nolint: exhaustivestruct
	msg := &dns.Msg{}

	// log.Printf("initialized worker using DNS server %s/udp\n", socket)

	return &checker{
		client: client,
		config: c,
		conn:   conn,
		msg:    msg,
	}, nil
}

//nolint: gochecknoinits
func init() {
	rand.Seed(time.Now().UnixNano())
}
