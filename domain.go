package tldcheck

import (
	"fmt"
)

type Domain struct {
	SLD string
	TLD
}

func (d Domain) Domain() string {
	return fmt.Sprintf("%s.%s", d.SLD, d.TLD.TLD)
}

func (d Domain) DomainRaw() string {
	return fmt.Sprintf("%s.%s", d.SLD, d.TLDRaw)
}

func (d Domain) String() string {
	if d.TLD.TLD != d.TLDRaw {
		return fmt.Sprintf("%s (%s)", d.Domain(), d.DomainRaw())
	}
	return d.Domain()
}
