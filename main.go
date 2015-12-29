package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

func main() {

	gen, err := newConfigFileReader("urls.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}

	processor := newDumpProcessor()
	for result := range processor.Process(gen.Start()) {
		fmt.Fprintf(os.Stdout, "url: %s\n", result)
	}

}

type Link struct {
	Name string
	URL  *url.URL
}

type UrlGenerator interface {
	Start() <-chan string
}

type UrlProcessor interface {
	Process(in <-chan string) <-chan string
}

type UrlFilter interface {
	Pass(Link) bool
}

type dumpProcessor struct {
	Client *http.Client
	out    chan string
}

func newDumpProcessor() UrlProcessor {
	return &dumpProcessor{Client: &http.Client{}, out: make(chan string)}
}

func (d *dumpProcessor) Process(in <-chan string) <-chan string {
	go func() {
		for r := range getURLs(getResponse(makeRequests(in), d.Client)) {
			d.out <- r
		}
		d.cleanup(nil)
	}()
	return d.out
}

func (d *dumpProcessor) cleanup(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	close(d.out)
}

func makeRequests(in <-chan string) <-chan *http.Request {
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

func getResponse(in <-chan *http.Request, client *http.Client) <-chan *http.Response {
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

func getURLs(in <-chan *http.Response) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

	responseLoop:
		for r := range in {
			z := html.NewTokenizer(r.Body)

			var inLink bool
			var link *url.URL
			var name string

			for {
				tt := z.Next()
				switch tt {
				case html.ErrorToken:
					log.Printf("HTML ERROR: %v\n", z.Err())
					continue responseLoop
				case html.TextToken:
					if inLink {
						name = string(z.Text())
					}
				case html.StartTagToken:
					tn, a := z.TagName()
					if len(tn) == 1 && tn[0] == 'a' && a {
						for {
							attr, val, more := z.TagAttr()
							if len(attr) == 4 && string(attr) == "href" {
								u, err := url.Parse(string(val))
								if err == nil {
									inLink = true
									link = u
								} else {
									log.Printf("error: %v", err)
								}
							}
							if !more || inLink {
								break
							}
						}
					}
				case html.EndTagToken:
					tn, _ := z.TagName()
					if len(tn) == 1 && tn[0] == 'a' {
						if inLink {
							out <- fmt.Sprintf("%q:%q", name, link)
							name = ""
							link = nil
							inLink = false
						}
					}
				case html.SelfClosingTagToken:
				case html.CommentToken:
				case html.DoctypeToken:
				}
			}
		}
	}()
	return out
}
