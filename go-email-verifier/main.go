package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMX, hasSPF, spfRecords, hasDMARC, DMARCRecords\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error: could not read from input: %v\n", err)
	}
}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecords string

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error looking up MX: %v", err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error looking up  TXT: %v", err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecord, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error looking up dmarc: %v", err)
	}

	for _, record := range dmarcRecord {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecords = record
			break
		}
	}

	fmt.Printf("%v,%v,%v,%v,%v,%v", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecords)
}
