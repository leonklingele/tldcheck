package tldcheck

import (
	"log" //nolint:depguard // TODO: Switch to log/slog
	"sync"

	"github.com/miekg/dns"
)

func Check(sld string, tlds []TLD, workers int) (ResultChan, error) {
	config := &dns.ClientConfig{
		Servers:  DNSServers,
		Port:     DNSPort,
		Timeout:  int(DNSTimeout.Seconds()),
		Attempts: DNSAttempts,
	}

	var wg sync.WaitGroup
	wch := make(chan workItem)
	rch := make(chan Result)

	for range workers {
		c, err := newChecker(config)
		if err != nil {
			return nil, err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			for wi := range wch {
				for {
					res, err := c.Check(wi)
					if err != nil {
						if wi.attempts <= 0 {
							log.Printf("failed to check item %s: %v\n", wi.String(), err)
							break
						}
						wi.attempts--
						continue
					}
					rch <- *res

					// Success
					break
				}
			}
		}()
	}

	go func() {
		for _, tld := range tlds {
			wi := workItem{
				domain: Domain{
					SLD: sld,
					TLD: tld,
				},
				attempts: config.Attempts,
			}
			wch <- wi
		}

		close(wch)
		wg.Wait()
		close(rch)
	}()

	return rch, nil
}
