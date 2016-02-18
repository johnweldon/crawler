package proc

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"

	"github.com/johnweldon/crawler/data"
	"github.com/johnweldon/crawler/util"
)

func NewExtractURLsProcessor(filter data.LinkFilter) URLProcessor {
	return &extractURLsProcessor{
		coreURLsProcessor: coreURLsProcessor{
			client: util.WebClient(),
			out:    make(chan string),
		},
		filter: filter,
	}
}

type extractURLsProcessor struct {
	coreURLsProcessor
	filter data.LinkFilter
}

func (d *extractURLsProcessor) Process(in <-chan string) <-chan string {
	go func() {
		for r := range d.GetURLs(in) {
			if d.filter.Pass(r) {
				d.out <- r.URL()
			}
		}
		d.cleanup(nil)
	}()
	return d.out
}

func (d *extractURLsProcessor) GetURLs(in <-chan string) <-chan *data.Link {
	return d.getLinks(d.Execute(in))
}

func (d *extractURLsProcessor) getLinks(in <-chan *http.Response) <-chan *data.Link {
	out := make(chan *data.Link)

	go func() {
		defer close(out)

	responseLoop:
		for r := range in {
			z := html.NewTokenizer(r.Body)

			var inLink bool
			var link *url.URL
			var name string
			var path data.StringStack

			for {
				tt := z.Next()
				switch tt {
				case html.ErrorToken:
					fmt.Fprintf(os.Stderr, "HTML ERROR: %v\n", z.Err())
					continue responseLoop
				case html.TextToken:
					if inLink {
						name = string(z.Text())
					}
				case html.StartTagToken:
					tn, a := z.TagName()
					path.Push(string(tn))
					if len(tn) == 1 && tn[0] == 'a' && a {
						for {
							attr, val, more := z.TagAttr()
							if len(attr) == 4 && string(attr) == "href" {
								u, err := url.Parse(string(val))
								if err == nil {
									inLink = true
									link = u
								} else {
									fmt.Fprintf(os.Stderr, "error: %v\n", err)
								}
							}
							if !more || inLink {
								break
							}
						}
					}
				case html.EndTagToken:
					tn, _ := z.TagName()
					if string(tn) == path.Peek() {
						path.Pop()
					}
					if len(tn) == 1 && tn[0] == 'a' {
						if inLink {
							out <- data.NewLink(name, link, path.String())
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
