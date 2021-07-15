package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/leonklingele/tldcheck"
	"github.com/logrusorgru/aurora/v3"
)

func run() error {
	dn := flag.String("name", "", "domain name to check for (required)")
	workers := flag.Int("workers", runtime.NumCPU(), "number of concurrent workers")
	flag.Parse()

	if len(*dn) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	tlds, err := tldcheck.AllTLDs()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	c, err := tldcheck.Check(*dn, tlds, *workers)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for res := range c {
		val := aurora.Red("not available")
		if res.Available {
			val = aurora.Green("available")
		}
		fmt.Printf("%s: %s\n", val, res.Domain.String()) //nolint: forbidigo
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
