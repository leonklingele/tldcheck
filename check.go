package tldcheck

import (
	"log"
	"sync"

	"github.com/miekg/dns"
)

func Check(sld string, tlds []TLD, workers int) (ResultChan, error) {
	//nolint: exhaustivestruct
	config := &dns.ClientConfig{
		Servers:  DNSServers,
		Port:     DNSPort,
		Timeout:  int(DNSTimeout.Seconds()),
		Attempts: DNSAttempts,
	}

	var wg sync.WaitGroup
	wch := make(chan workItem)
	rch := make(chan Result)

	for i := 0; i < workers; i++ {
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
							break //nolint: nlreturn
						}
						wi.attempts--
						continue //nolint: nlreturn
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
