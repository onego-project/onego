package resources

import (
	"fmt"
	"github.com/beevik/etree"
	"strconv"
)

// Resource structure contains XML data and main methods for open nebula resources
type Resource struct {
	XMLData *etree.Element
}

const invalidCode = -1

// constants for value conversion in resources methods
const base10 = 10
const bitSize64 = 64

// CreateResource constructs Resource by tag and ID
func CreateResource(tag string, id int) *Resource {
	doc := etree.NewDocument()

	el := doc.CreateElement(tag).CreateElement("ID")
	el.SetText(fmt.Sprintf("%d", id))

	return &Resource{doc.Root()}
}

// Attribute gets resource attribute founded on the path
func (r *Resource) Attribute(path string) (string, error) {
	if r.XMLData == nil {
		return "", fmt.Errorf("no xml data, unable to get %s", path)
	}

	element := r.XMLData.FindElement(path)
	if element == nil {
		return "", fmt.Errorf("unable to find %s", path)
	}
	return element.Text(), nil
}

func (r *Resource) intAttribute(path string) (int, error) {
	attribute, err := r.Attribute(path)
	if err != nil {
		return invalidCode, err
	}
	i, err := strconv.Atoi(attribute)
	if err != nil {
		return invalidCode, err
	}
	return i, nil
}

// ID gets resource id
func (r *Resource) ID() (int, error) {
	return r.intAttribute("ID")
}

// Name gets resource name
func (r *Resource) Name() (string, error) {
	return r.Attribute("NAME")
}

func intToBool(i int) bool {
	return i == 1
}
