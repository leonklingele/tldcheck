package main

import (
	"flag" //nolint:depguard // We only allow to import the flag package in here
	"fmt"
	"log" //nolint:depguard // TODO: Switch to log/slog
	"os"
	"runtime"

	"github.com/leonklingele/tldcheck"

	"github.com/logrusorgru/aurora/v3"
)

func run() error {
	dn := flag.String("name", "", "domain name to check for (required)")
	workers := flag.Int("workers", runtime.NumCPU(), "number of concurrent workers")
	flag.Parse()

	if *dn == "" {
		flag.Usage()
		os.Exit(1) //nolint:revive // Fine to do in the main package
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

		//nolint:errcheck,forbidigo // We explicitly want to print to stdout
		_, _ = fmt.Printf("%s: %s\n", val, res.Domain.String())
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
