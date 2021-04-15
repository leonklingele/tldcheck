package tldcheck

import (
	_ "embed"
	"fmt"
	"strings"

	"golang.org/x/net/idna"
)

//go:embed tlds-alpha-by-domain.txt
var rawTLDList string //nolint: gochecknoglobals

type TLD struct {
	TLD    string
	TLDRaw string
}

func AllTLDs() ([]TLD, error) {
	rawTLDs := strings.Split(rawTLDList, "\n")
	tlds := make([]TLD, 0, len(rawTLDs))

	for _, tldRaw := range rawTLDs {
		if len(tldRaw) == 0 || strings.HasPrefix(tldRaw, "#") {
			continue
		}

		tldRaw = strings.ToLower(tldRaw)
		tld := tldRaw

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
