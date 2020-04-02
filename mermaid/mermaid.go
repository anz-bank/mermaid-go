package mermaid

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

var flagPort = flag.Int("port", 8544, "port")

type diagram struct {
	html string
	svg  string
}

func Execute(mermaidCode ...string) string {
	go Server(fmt.Sprintf(":%d", *flagPort), LoadTemplate(mermaidCode...))
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	d := &diagram{}
	if err := chromedp.Run(ctx, d.browserTasks(fmt.Sprintf("http://localhost:%d", *flagPort))); err != nil {
		log.Fatal(err)
	}

	r := strings.NewReader(d.html)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		panic(err)
	}

	doc.Find("#mermaid").Each(func(i int, s *goquery.Selection) {
		inside_html, err := s.Html()
		if err != nil {
			panic(err)
		}
		d.svg = inside_html
	})
	return d.svg
}

func (s *diagram) browserTasks(host string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(host),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate("").Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			s.html, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	}
}
