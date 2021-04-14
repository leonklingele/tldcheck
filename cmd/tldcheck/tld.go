package main

import (
	_ "embed"
	"fmt"
	"strings"

	"golang.org/x/net/idna"
)

//go:embed tlds-alpha-by-domain.txt
var rawTLDList string //nolint: gochecknoglobals

func tlds() (map[string]string, error) {
	tlds := make(map[string]string)

	rawTLDs := strings.Split(rawTLDList, "\n")
	for _, tld := range rawTLDs {
		if len(tld) == 0 || strings.HasPrefix(tld, "#") {
			continue
		}

		tld = strings.ToLower(tld)

		if strings.HasPrefix(tld, "xn--") {
			// Punycode TLD
			u, err := idna.ToUnicode(tld)
			if err != nil {
				return nil, fmt.Errorf("failed to convert un-puny-fy TLD %q: %w", tld, err)
			}
			tlds[tld] = u
		} else {
			tlds[tld] = tld
		}
	}

	return tlds, nil
}
