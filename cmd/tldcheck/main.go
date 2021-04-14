package main

import (
	"flag"
	"log"
	"os"
	"runtime"
)

func run() error {
	dn := flag.String("name", "", "domain name to check for")
	workers := flag.Int("workers", runtime.NumCPU(), "number of concurrent workers")
	flag.Parse()

	if len(*dn) == 0 {
		log.Println("no domain name given")
		os.Exit(1)
	}

	tlds, err := tlds()
	if err != nil {
		return err
	}

	return check(*dn, tlds, *workers)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
