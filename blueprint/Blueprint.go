package blueprint

import (
	"bytes"

	"github.com/onego-project/onego/errors"

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

// CreateBlueprint creates Blueprint with root
func CreateBlueprint(rootElement string) *Blueprint {
	doc := etree.NewDocument()
	doc.CreateElement(rootElement)

	return &Blueprint{XMLData: doc}
}

// Render method to render blueprint values to string
func (bp *Blueprint) Render() (string, error) {
	if bp.XMLData == nil {
		return "", errors.ErrBlueprintXMLEmpty
	}

	var buffer bytes.Buffer
	if _, err := bp.XMLData.WriteTo(&buffer); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// SetElement sets element to blueprint with tag and value
func (bp *Blueprint) SetElement(tag, value string) {
	element := bp.XMLData.Root().FindElement(tag)
	if element == nil {
		element = bp.XMLData.Root().CreateElement(tag)
	}
	element.SetText(value)
}
