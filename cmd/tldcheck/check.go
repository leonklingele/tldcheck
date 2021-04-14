package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/logrusorgru/aurora/v3"
	"github.com/miekg/dns"
)

const (
	dnsTimeout  = 8 * time.Second
	dnsAttempts = 8
)

var (
	//nolint: gochecknoglobals
	dnsServers = []string{
		"1.1.1.1",
		"8.8.4.4",
		"8.8.8.8",
		"9.9.9.9",
	}
	//nolint: gochecknoglobals
	dnsPort = "53"
)

type availabilityInfo struct {
	available bool
}

type workItem struct {
	dn       string
	rawTLD   string
	tld      string
	attempts int
}

func (wi *workItem) DomainRaw() string {
	return fmt.Sprintf("%s.%s", wi.dn, wi.rawTLD)
}

func (wi *workItem) Domain() string {
	return fmt.Sprintf("%s.%s", wi.dn, wi.tld)
}

func (wi *workItem) String() string {
	if wi.rawTLD != wi.tld {
		return fmt.Sprintf("%s (%s)", wi.Domain(), wi.DomainRaw())
	}
	return wi.Domain() //nolint: nlreturn
}

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

func (c *checker) Check(wi *workItem) (*availabilityInfo, error) {
	r, err := c.dnsQuery(wi.DomainRaw())
	if err != nil {
		return nil, err
	}

	var available bool
	if r.Rcode == dns.RcodeNameError {
		available = true
	}

	return &availabilityInfo{
		available: available,
	}, nil
}

func newChecker(c *dns.ClientConfig) (*checker, error) {
	//nolint: exhaustivestruct
	client := &dns.Client{
		Timeout: dnsTimeout,
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

	// log.Printf("initialized worker using DNS server %s:%s/udp\n", server, port)

	return &checker{
		client: client,
		config: c,
		conn:   conn,
		msg:    msg,
	}, nil
}

func check(dn string, tlds map[string]string, workers int) error {
	//nolint: exhaustivestruct
	config := &dns.ClientConfig{
		Servers:  dnsServers,
		Port:     dnsPort,
		Timeout:  int(dnsTimeout.Seconds()),
		Attempts: dnsAttempts,
	}

	var wg sync.WaitGroup
	wch := make(chan *workItem)

	for i := 0; i < workers; i++ {
		c, err := newChecker(config)
		if err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			for wi := range wch {
				for {
					ai, err := c.Check(wi)
					if err != nil {
						if wi.attempts <= 0 {
							log.Printf("failed to check item %s: %v\n", wi.String(), err)
							break //nolint: nlreturn
						}
						wi.attempts--
						continue //nolint: nlreturn
					}

					val := aurora.Red("not available")
					if ai.available {
						val = aurora.Green("available")
					}
					fmt.Printf("%s: %s\n", val, wi.String()) //nolint: forbidigo

					// Success
					break
				}
			}
		}()
	}

	for rawTLD, tld := range tlds {
		wi := &workItem{
			dn:       dn,
			rawTLD:   rawTLD,
			tld:      tld,
			attempts: config.Attempts,
		}
		wch <- wi
	}

	close(wch)
	wg.Wait()

	return nil
}

//nolint: gochecknoinits
func init() {
	rand.Seed(time.Now().UnixNano())
}
