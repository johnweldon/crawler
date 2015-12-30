package proc

import (
	"fmt"
	"net/http"
	"os"
)

type coreURLsProcessor struct {
	client *http.Client
	out    chan string
}

func (d *coreURLsProcessor) cleanup(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	close(d.out)
}

func (d *coreURLsProcessor) makeRequests(in <-chan string) <-chan *http.Request {
	out := make(chan *http.Request)
	go func() {
		for u := range in {
			r, err := http.NewRequest("GET", u, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: url %q -> %v\n", u, err)
			} else {
				out <- r
			}
		}
		close(out)
	}()
	return out
}

func (d *coreURLsProcessor) getResponse(in <-chan *http.Request, client *http.Client) <-chan *http.Response {
	out := make(chan *http.Response)
	go func() {
		for r := range in {
			resp, err := client.Do(r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			} else {
				out <- resp
			}
		}
		close(out)
	}()
	return out
}

func (d *coreURLsProcessor) Execute(in <-chan string) <-chan *http.Response {
	return d.getResponse(d.makeRequests(in), d.client)
}
