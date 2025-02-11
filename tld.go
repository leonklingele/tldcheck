package tldcheck

import (
	"fmt"
	"strings"

	"golang.org/x/net/idna"

	_ "embed"
)

//go:embed tlds-alpha-by-domain.txt
var rawTLDList string

type TLD struct {
	TLD    string
	TLDRaw string
}

func AllTLDs() ([]TLD, error) {
	rawTLDs := strings.Split(rawTLDList, "\n")
	tlds := make([]TLD, 0, len(rawTLDs))

	for _, tldRaw := range rawTLDs {
		if tldRaw == "" || strings.HasPrefix(tldRaw, "#") {
			continue
		}

		tldRaw = strings.ToLower(tldRaw)
		tld := tldRaw //nolint:copyloopvar // False positive as used as a reference to the previous version

		if strings.HasPrefix(tldRaw, "xn--") {
			// Punycode TLD
			u, err := idna.ToUnicode(tldRaw)
			if err != nil {
				return nil, fmt.Errorf("failed to un-puny-fy TLD %q: %w", tldRaw, err)
			}

			tld = u
		}

		tlds = append(tlds, TLD{
			TLD:    tld,
			TLDRaw: tldRaw,
		})
	}

	return tlds, nil
}
