package main

import (
	"fmt"
	"os"

	"github.com/johnweldon/crawler/crawl"
)

func main() {
	if source := os.Getenv("URL_FILE"); source == "" {
		fmt.Fprintf(os.Stderr, "no URL_FILE specified\n")
	} else {
		crawl.GetAllLinksInTable(source)
	}

	if source := os.Getenv("TABLE_FILE"); source == "" {
		fmt.Fprintf(os.Stderr, "no TABLE_FILE specified\n")
	} else {
		crawl.GetTables(source)
	}
}
