package blueprint

import (
	"bytes"
	"fmt"
	"github.com/beevik/etree"
)

// Blueprint structure
type Blueprint struct {
	XMLData *etree.Document
}

// Interface of Blueprint
type Interface interface {
	Render() (string, error)
}

// CreateBlueprint creates empty Blueprint
func CreateBlueprint() *Blueprint {
	return &Blueprint{XMLData: etree.NewDocument()}
}

// Render method to render blueprint values to string
func (bp *Blueprint) Render() (string, error) {
	if bp.XMLData == nil {
		return "", fmt.Errorf("blueprint XML data is empty")
	}

	var buffer bytes.Buffer
	if _, err := bp.XMLData.WriteTo(&buffer); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// SetElement sets element to blueprint with tag and value
func (bp *Blueprint) SetElement(tag, value string) {
	template := bp.XMLData.FindElement("TEMPLATE")
	if template == nil {
		template = bp.XMLData.CreateElement("TEMPLATE")
	}

	element := template.FindElement(tag)
	if element == nil {
		element = template.CreateElement(tag)
	}
	element.SetText(value)
}
