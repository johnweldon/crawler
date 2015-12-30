package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/johnweldon/crawler/data"
	"github.com/johnweldon/crawler/gen"
	"github.com/johnweldon/crawler/proc"
)

func main() {
	urls := os.Getenv("URL_FILE")
	if urls == "" {
		fmt.Fprintf(os.Stderr, "no URL_FILE specified")
		return
	}

	gen, err := gen.NewConfigFileReader(urls)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}

	urls := map[string]interface{}{}
	processor := proc.NewExtractURLsProcessor(data.NewLinkInTableFilter())
	for result := range processor.Process(gen.Start()) {
		urls[result] = nil
	}

	client := &http.Client{}
	for k, _ := range urls {
		dump(client, k)
	}
}

func dump(c *http.Client, link string) {
	u, err := url.Parse(link)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error with link: %q, %v", link, err)
		return
	}

	p := spaceMap(filepath.Base(u.Path))

	f, err := os.Create(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating file: %q, %v", p, err)
		return
	}
	defer f.Close()

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error downloading link: %q, %v", link, err)
		return
	}

	if a := os.Getenv("USER_AGENT"); a != "" {
		req.Header.Set("User-Agent", a)
	}

	r, err := c.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error downloading link: %q, %v", link, err)
		return
	}
	defer r.Body.Close()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error saving link: %q, %v", link, err)
		return
	}

	return
}

// adapted from http://stackoverflow.com/a/32081891/102371
func spaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '_'
		}
		return r
	}, str)
}
