package main

import (
	"log"
	"testing"

	"github.com/joshcarp/mermaid-go/mermaid"
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
	result := mermaid.Execute(string(file))

	outfile, err := fs.Create(output)
	if err != nil {
		log.Fatal("Error: creating output file")
	}
	outfile.Write([]byte(result))
}
