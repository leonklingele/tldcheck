package tldcheck

import (
	"fmt"
	rand "math/rand/v2"
	"time"

	"github.com/miekg/dns"
)

//nolint:gochecknoglobals // Nice to have as a global
var (
	DNSServers = []string{
		"1.1.1.1",
		"8.8.4.4",
		"8.8.8.8",
		"9.9.9.9",
	}
	DNSPort = "53"

	DNSTimeout  = 8 * time.Second
	DNSAttempts = 8
)

type checker struct {
	client *dns.Client
	conn   *dns.Conn
	msg    *dns.Msg
}

func (c *checker) dnsQuery(domain string) (*dns.Msg, error) {
	c.msg.SetQuestion(dns.Fqdn(domain), dns.TypeNS)

	r, _, err := c.client.ExchangeWithConn(c.msg, c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", domain, err)
	}
	return r, nil
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
	client := &dns.Client{
		Timeout: DNSTimeout,
	}
	i := rand.IntN(len(c.Servers)) //nolint:gosec // Must not be a crypto-rand index
	server := c.Servers[i]
	port := c.Port
	socket := fmt.Sprintf("%s:%s", server, port)
	conn, err := client.Dial(socket)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", socket, err)
	}
	msg := &dns.Msg{}

	return &checker{
		client: client,
		conn:   conn,
		msg:    msg,
	}, nil
}
