package blueprint

import (
	"bytes"
	"github.com/beevik/etree"
	"strings"
)

// Blueprint structure
type Blueprint struct {
	Values map[string]Node
	Root   string
}

// Node struct
type Node struct {
	Value  string
	Values map[string]Node
}

// Render method to render blueprint values to string
func (bp Blueprint) Render() string {
	doc := etree.NewDocument()
	root := doc.CreateElement(strings.ToUpper(bp.Root))
	for k, v := range bp.Values {
		e := root.CreateElement(strings.ToUpper(k))
		v.renderXML(e)
	}
	var buffer bytes.Buffer
	if _, err := doc.WriteTo(&buffer); err != nil {
		return ""
	}

	return buffer.String()
}

func (bn Node) renderXML(e *etree.Element) {
	if bn.Value != "" {
		e.SetText(bn.Value)
	} else {
		for k, v := range bn.Values {
			f := e.CreateElement(strings.ToUpper(k))
			v.renderXML(f)
		}
	}
}

// Interface of Blueprint
type Interface interface {
	Render() string
}
