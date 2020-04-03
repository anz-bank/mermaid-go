package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/joshcarp/mermaid-go/mermaid"
)

func main() {
	var output string
	flag.StringVar(&output, "o", "", "Output file of the svg")
	flag.Parse()
	filename := flag.Arg(0)
	if filename == "" {
		log.Fatal("Error: no filename specified")
	}
	if output == "" {
		output = strings.TrimSuffix(filename, filepath.Ext(filename)) + ".svg"
	}
	fs := afero.NewOsFs()
	file, err := afero.ReadFile(fs, filename)
	if err != nil {
		log.Fatal("Error: reading input file")
	}
	result := mermaid.Execute(string(file))

	outfile, err := fs.Create(output)
	if err != nil {
		log.Fatal("Error: creating output file")
	}
	outfile.Write([]byte(result[0]))

}
