# mermaid-go

[mermaid-js](https://github.com/mermaid-js/mermaid) executable written in Go.

Uses [chromedp](https://github.com/chromedp/chromedp) to run mermaid js.

## Use

### CLI 
#### Installation
`go get -u github.com/anz-bank/mermaid-go`
#### Execute
`mermaid-go <input.mmdc> -o <output.svg>`

See [demo](demo)

### As a go package

`go get -u github.com/anz-bank/mermaid-go/mermaid`


```
package main

import (
	"io/ioutil"

	"github.com/anz-bank/mermaid-go/mermaid"
)

func main() {
	str := `
graph TD;
A-->B;
A-->C;
B-->D;
C-->D;
	`
	g := mermaid.Init()
	svg := g.Execute(str)
	if err := ioutil.WriteFile("mermaid.svg", []byte(svg), 0644); err != nil {
		panic(err)
	}
}

```
