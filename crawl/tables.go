package crawl

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"text/template"

	"github.com/johnweldon/crawler/data"
	"github.com/johnweldon/crawler/gen"
	"github.com/johnweldon/crawler/proc"
)

func GetTables(urlfile string) {
	gen, err := gen.NewConfigFileReader(urlfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}

	processor := proc.NewExtractTablesProcessor()

	for result := range processor.Process(gen.Start()) {
		saveTable(result)
	}
}

func saveTable(table *data.Table) {
	u, err := url.Parse(table.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error with name: %q, %v\n", table.Name, err)
		return
	}

	p := spaceMap(filepath.Base(u.Path))

	f, err := os.Create(fmt.Sprintf("%s_%d.html", p, table.Ordinal))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating file: %q, %v\n", p, err)
		return
	}
	defer f.Close()

	tmpl := template.Must(template.New("out").Parse(`<html>
  <head>
	  <title>{{ .Name }} {{ .Ordinal }}</title>
	</head>
	<body>
	  <table>
		  <thead>
			  <tr>{{ range .Header }}
				  <th>{{ . }}</th>{{ end }}
				</tr>
			</thead>
			<tbody>{{ range .Rows }}
			  <tr>{{ range . }}
				<td>{{ . }}</td>{{ end }}
				</tr>
{{ end }}
			</tbody>
		</table>
	</body>
</html>
`))

	tmpl.Execute(f, table)
}
