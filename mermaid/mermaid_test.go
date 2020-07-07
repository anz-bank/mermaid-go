package mermaid

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const sampleClassDiagram = `classDiagram
Animal <|-- Duck
Animal <|-- Fish
Animal <|-- Zebra
Animal : +int age
Animal : +String gender
Animal: +isMammal()
Animal: +mate()
class Duck{
	+String beakColor
	+swim()
	+quack()
}
class Fish{
	-int sizeInFeet
	-canEat()
}
class Zebra{
	+bool is_wild
	+run()
}
				

`

const sampleFlowChart = `graph TD
A[Christmas] -->|Get money| B(Go shopping)
B --> C{Let me think}
C -->|One| D[Laptop]
C -->|Two| E[iPhone]
C -->|Three| F[fa:fa-car Car]
				
`

const sampleSequenceDiagram = `sequenceDiagram
Alice->>+John: Hello John, how are you?
Alice->>+John: John, can you hear me?
John-->>-Alice: Hi Alice, I can hear you!
John-->>-Alice: I feel great!
				
`

const sampleErDiagram = `erDiagram
CUSTOMER }|..|{ DELIVERY-ADDRESS : has
CUSTOMER ||--o{ ORDER : places
CUSTOMER ||--o{ INVOICE : "liable for"
DELIVERY-ADDRESS ||--o{ ORDER : receives
INVOICE ||--|{ ORDER : covers
ORDER ||--|{ ORDER-ITEM : includes
PRODUCT-CATEGORY ||--|{ PRODUCT : contains
PRODUCT ||--o{ ORDER-ITEM : "ordered in"
`

const classDiagramWithList = `classDiagram
class Server_Response {
	MegaDatabase_Money query
	MegaDatabase_Empty balance
}
class Server_Request {
	List<Server_Response> query
}
Server_Request <-- Server_Response

`

func TestExecuteClassDiagram(t *testing.T) {
	svgResult := Execute(sampleClassDiagram)
	assert.NoError(t, isValidMermaidSVG(svgResult))
	err := ioutil.WriteFile("./test/sampleClassDiagram.svg", []byte(svgResult), 0644)
	assert.NoError(t, err)
}

func TestExecuteFlowChart(t *testing.T) {
	svgResult := Execute(sampleFlowChart)
	assert.NoError(t, isValidMermaidSVG(svgResult))
	err := ioutil.WriteFile("./test/sampleFlowChart.svg", []byte(svgResult), 0644)
	assert.NoError(t, err)
}

func TestExecuteSequenceDiagram(t *testing.T) {
	svgResult := Execute(sampleSequenceDiagram)
	assert.NoError(t, isValidMermaidSVG(svgResult))
	err := ioutil.WriteFile("./test/sampleSequenceDiagram.svg", []byte(svgResult), 0644)
	assert.NoError(t, err)
}

func TestExecuteErDiagram(t *testing.T) {
	svgResult := Execute(sampleErDiagram)
	assert.NoError(t, isValidMermaidSVG(svgResult))
	err := ioutil.WriteFile("./test/sampleErDiagram.svg", []byte(svgResult), 0644)
	assert.NoError(t, err)
}

func TestExecuteClassDiagramWithList(t *testing.T) {
	svgResult := Execute(classDiagramWithList)
	assert.NoError(t, isValidMermaidSVG(svgResult))
	err := ioutil.WriteFile("./test/classDiagramWithList.svg", []byte(svgResult), 0644)
	assert.NoError(t, err)
}

// isValidMermaidSVG checks the validity of the generated svg
// For now it only checks for the presence of <svg> and <defs> tags
func isValidMermaidSVG(svgResult string) error {
	switch {
	case svgResult == "":
		return errors.New("SVG is an empty string")
	case !strings.Contains(svgResult, "<defs>"):
		return errors.New("SVG doesn't contain <defs> element")
	case xml.Unmarshal([]byte(svgResult), new(interface{})) != nil:
		return errors.New("SVG is invalid xml")
	default:
		return nil
	}
}
