package proc

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"

	"github.com/johnweldon/crawler/data"
	"github.com/johnweldon/crawler/util"
)

type ExtractTablesProcessor struct {
	coreURLsProcessor
}

func NewExtractTablesProcessor() *ExtractTablesProcessor {
	return &ExtractTablesProcessor{
		coreURLsProcessor: coreURLsProcessor{
			client: util.WebClient(),
			out:    make(chan string),
		}}
}

func (p *ExtractTablesProcessor) Process(in <-chan string) <-chan *data.Table {
	out := make(chan *data.Table)
	go func() {
		defer close(out)
		for table := range p.getTables(p.Execute(in)) {
			out <- table
		}
	}()
	return out
}

func (p *ExtractTablesProcessor) getTables(in <-chan *http.Response) <-chan *data.Table {
	out := make(chan *data.Table)

	go func() {
		defer close(out)

	responseLoop:
		for r := range in {
			z := html.NewTokenizer(r.Body)

			var doc = r.Request.URL.String()
			var count = 0
			var inTable bool
			var inTD bool
			var inTH bool
			var name string
			var row []string
			var head []string
			var table = new(data.Table)

			for {
				tt := z.Next()
				switch tt {
				case html.ErrorToken:
					fmt.Fprintf(os.Stderr, "HTML ERROR: %v\n", z.Err())
					continue responseLoop
				case html.TextToken:
					if inTD || inTH {
						name = fmt.Sprintf("%s %s", name, strings.TrimSpace(string(z.Text())))
					}
				case html.StartTagToken:
					tn, _ := z.TagName()
					switch string(tn) {
					case "table":
						// TODO: handle nested tables, maybe with new parser
						if !inTable {
							table = &data.Table{Name: doc, Ordinal: count}
							inTable = true
						}
					case "td":
						name = ""
						inTD = true
					case "th":
						name = ""
						inTH = true
					case "tr":
						row = []string{}
						head = []string{}
					}
				case html.EndTagToken:
					tn, _ := z.TagName()
					switch string(tn) {
					case "table":
						if table != nil && len(table.Rows) > 0 {
							out <- table
						}
						table = nil
						count++
						inTable = false
					case "td":
						row = append(row, strings.TrimSpace(name))
						inTD = false
					case "th":
						head = append(head, name)
						inTH = false
					case "tr":
						if !inTable {
							table = &data.Table{Name: doc, Ordinal: count}
							inTable = true
						}
						if len(head) > 0 {
							table.Header = head
						}
						if len(row) > 0 {
							table.Rows = append(table.Rows, row)
						}
					}
				case html.SelfClosingTagToken:
					tn, _ := z.TagName()
					switch string(tn) {
					case "table":
						inTable = false
					case "td":
						inTD = false
					case "th":
						inTH = false
					}
				case html.CommentToken:
				case html.DoctypeToken:
				}
			}
		}
	}()
	return out
}
