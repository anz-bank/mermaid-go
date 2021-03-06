package main

import (
	"log"
	"testing"

	"github.com/anz-bank/mermaid-go/mermaid"
	"github.com/spf13/afero"
)

func TestY(t *testing.T) {
	output := "output.svg"
	filename := "demo/flowchart.mmd"
	fs := afero.NewOsFs()
	file, err := afero.ReadFile(fs, filename)
	fs = afero.NewMemMapFs()
	if err != nil {
		log.Fatal("Error: reading input file")
	}
	g := mermaid.Init()
	result := g.Execute(string(file))

	outfile, err := fs.Create(output)
	if err != nil {
		log.Fatal("Error: creating output file")
	}
	_, err = outfile.Write([]byte(result))
	if err != nil {
		log.Fatal("Error: saving output file")
	}
}
