package tldcheck

type workItem struct {
	domain   Domain
	attempts int
}

func (wi workItem) String() string {
	return wi.domain.String()
}
