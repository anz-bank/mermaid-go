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
var str string

type svg struct {
	content string
}

func Execute(ctx context.Context) {
	s := &svg{}
	// run task list
	err := chromedp.Run(ctx, s.browserTasks(fmt.Sprintf("http://localhost:%d", *flagPort)))
	if err != nil {
		log.Fatal(err)
	}

	r := strings.NewReader(str)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		panic(err)
	}
	doc.Find("#mermaid").Each(func(i int, s *goquery.Selection) {
		inside_html, _ := s.Html() //underscore is an error
		fmt.Printf("Review %d: %s\n", i, inside_html)
	})
}

func (s *svg) browserTasks(host string) chromedp.Tasks {
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
			str, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	}
}
