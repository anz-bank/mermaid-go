package mermaid

import (
	"bytes"
	"context"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/chromedp/chromedp"
)

// Execute evaluates mermaid code and returns svg string slices.
func Execute(mermaidCode string) string {
	return EvaluateAndSelectHTML(LoadTemplate(mermaidCode), "svg")
}

func writeHTML(content string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, strings.TrimSpace(content))
	})
}

// Execute evaluates raw html (with javascript embedded) and returns the processed HTML.
func EvaluateAndSelectHTML(rawHTML, selector string) string {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ts := httptest.NewServer(writeHTML(rawHTML))
	defer ts.Close()
	var processedHTML string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(ts.URL),
		chromedp.OuterHTML(selector, &processedHTML),
	); err != nil {
		panic(err)
	}
	return processedHTML
}

const indexHTML = `<html>
<body>
<script src="https://cdn.jsdelivr.net/npm/mermaid@8.4.0/dist/mermaid.min.js"></script>
<script>mermaid.initialize({startOnLoad:true});</script>
<div class="mermaid">
{{.}}
</div>
</body>
</html>`

// Load template returns a mermaid html page with the input mermaid code embedded.
func LoadTemplate(file string) string {
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
