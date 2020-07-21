// +build chrome

// These tests call mermaid-go concurrently in parallal, to help simulate the behaviour when
// mermaid-go is used to generate large numbers of diagrams.
package mermaid

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestConcurrentSharedChrome is an example of the right way to call mermaid-go. Runs all operations in the same browser instance
func TestConcurrentSharedChrome(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup
	g := Init()
	start := time.Now()
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			g.Execute(sampleClassDiagram)
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	g.Close()
	elapsed := time.Since(start)
	t.Log(elapsed)
}

// TestConcurrentNewChrome is an example of the wrong way to call mermaid-go. It creates a new browser instance for each diagram generation call.
func TestConcurrentNewChrome(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			g := Init()
			defer g.Close()
			svgResult := g.Execute(sampleClassDiagram)
			assert.NoError(t, isValidMermaidSVG(svgResult))
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	elapsed := time.Since(start)
	t.Log(elapsed)
}

// TestSequentialSharedChrome illustrates the advantages of using a single chrome instance
// It takes 2.5s to generate 10 diagrams
func TestSequentialSharedChrome(t *testing.T) {
	t.Parallel()
	g := Init()
	start := time.Now()
	for i := 0; i < 10; i++ {
		g.Execute(sampleClassDiagram)
	}
	g.Close()
	elapsed := time.Since(start)
	t.Log(elapsed)
}

// TestSequentialNewChrome illustrates the performance overhead of creating a new chrome instance for each call.
// It takes 10s to generate 10 diagrams
func TestSequentialNewChrome(t *testing.T) {
	t.Parallel()
	start := time.Now()
	for i := 0; i < 10; i++ {
		g := Init()
		g.Execute(sampleClassDiagram)
		g.Close()
	}
	elapsed := time.Since(start)
	t.Log(elapsed)
}
