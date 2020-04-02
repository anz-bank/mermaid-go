package mermaid

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

const indexHTML = `<html>
<body>
<script src="https://cdn.jsdelivr.net/npm/mermaid@8.4.0/dist/mermaid.min.js"></script>
<script>mermaid.initialize({startOnLoad:true});</script>

Here is one mermaid diagram:
<div id="mermaid" class="mermaid">
{{range $file := .}}{{.}}{{end}}
</div>
</body>
</html>`

// Server is a simple HTTP server that serves a static html page.
func Server(addr, content string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(res, content)
	})
	return http.ListenAndServe(addr, mux)
}

func LoadTemplate(file ...string) string {
	newTemplate, err := template.New("").Parse(indexHTML)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	if err := newTemplate.Execute(&buf, file); err != nil {
		panic(err)
	}
	return buf.String()
}
