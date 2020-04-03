package mermaid

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
	"github.com/markbates/pkger"

	"github.com/chromedp/chromedp"
)

const mermaidJS = "github.com/joshcarp/mermaid-go:/resources/mermaid.min.js"

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

// Execute evaluates raw html (with javascript embedded) and returns the processed HTML.
func EvaluateAndSelectHTML(rawHTML, selector string) string {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(rawHTML))
	})
	r.Get("/mermaid.min.js", func(w http.ResponseWriter, r *http.Request) {
		w.Write(OpenFile(mermaidJS))
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

// OpenFile opens a statically encoded file from pkger.
func OpenFile(filename string) []byte {
	pkger.Include(filename)
	scriptFile, err := pkger.Open(filename)
	if err != nil {
		panic(err)
	}
	scriptBytes := make([]byte, 1113944)
	n, err := scriptFile.Read(scriptBytes)
	if err != nil {
		panic(err)
	}
	return scriptBytes[0:n]
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
