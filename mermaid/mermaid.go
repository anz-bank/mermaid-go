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

type Generator struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// // Init returns a new Generator instance which sets up a chrome browser instance to be used by all subsequent calls
func Init() *Generator {
	ctx, cancel := chromedp.NewContext(context.Background())
	g := &Generator{}
	g.ctx = ctx
	g.cancel = cancel
	if err := chromedp.Run(ctx); err != nil {
		panic(err)
	}
	return g
}

// Close cancels the root context and closes the chrome browser instance created by chromedp
func (g *Generator) Close() {
	g.cancel()
}

// Execute evaluates mermaid code and returns svg string slices.
func (g *Generator) Execute(mermaidCode string) string {
	return g.EvaluateAndSelectHTML(LoadTemplate(mermaidCode), "svg")
}

// Execute evaluates raw html (with javascript embedded) and returns the processed HTML.
func (g *Generator) EvaluateAndSelectHTML(rawHTML, selector string) string {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(rawHTML))
		if err != nil {
			panic(err)
		}
	})
	r.Get("/mermaid.min.js", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write((Decode64([]byte(mermaidjs64))))
		if err != nil {
			panic(err)
		}
	})
	ts := httptest.NewServer(r)
	defer ts.Close()
	var processedHTML string
	ctx, cancel := chromedp.NewContext(g.ctx)
	defer cancel()

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

func Decode64(src []byte) []byte {
	decoded := make([]byte, 2176000)
	n, err := base64.StdEncoding.Decode(decoded, src)
	if err != nil {
		panic(err)
	}
	return decoded[:n]
}
