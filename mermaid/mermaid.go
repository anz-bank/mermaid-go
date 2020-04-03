package mermaid

import (
	"bytes"
	"context"
	"encoding/base64"
	"html/template"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"

	"github.com/chromedp/chromedp"
)

const indexHTML = `<html>
<body>
<script src="mermaid.min.js"></script>
<script>mermaid.initialize({startOnLoad:true});</script>
<div class="mermaid">
{{.}}
</div>
</body>
</html>`

// Execute evaluates mermaid code and returns svg string slices.
func Execute(mermaidCode string) string {
	return EvaluateAndSelectHTML(LoadTemplate(mermaidCode), "svg")
}
func Decode64(src []byte) []byte {
	decoded := make([]byte, 2176000)
	n, err := base64.StdEncoding.Decode(decoded, src)
	if err != nil {
		panic(err)
	}
	return decoded[:n]
}

// Execute evaluates raw html (with javascript embedded) and returns the processed HTML.
func EvaluateAndSelectHTML(rawHTML, selector string) string {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(rawHTML))
	})
	r.Get("/mermaid.min.js", func(w http.ResponseWriter, r *http.Request) {
		w.Write((Decode64([]byte(mermaidjs64))))
	})
	ts := httptest.NewServer(r)
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

// Load template returns a mermaid html page with the input mermaid code embedded.
func LoadTemplate(mermaidSource string) string {
	newTemplate, err := template.New("").Parse(indexHTML)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	if err := newTemplate.Execute(&buf, mermaidSource); err != nil {
		panic(err)
	}
	return buf.String()
}
