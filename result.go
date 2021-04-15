package tldcheck

type Result struct {
	Domain    Domain
	Available bool
}

type ResultChan <-chan Result
