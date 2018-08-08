package resources

import (
	"fmt"
	"github.com/beevik/etree"
	"strconv"
)

// Host struct
type Host struct {
	XMLData *etree.Element
}

// CreateHost constructs Host
func CreateHost(id int) *Host {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0"`)

	el := doc.CreateElement("HOST").CreateElement("ID")
	el.SetText(fmt.Sprintf("%d", id))

	return &Host{doc.Root()}
}

// GetAttribute method
func (h Host) GetAttribute(path string) string {
	elements := h.XMLData.FindElements(path)
	if elements == nil {
		return ""
	}
	return elements[0].Text()
}

// GetID method
func (h Host) GetID() int {
	i, err := strconv.Atoi(h.GetAttribute("ID"))
	if err != nil {
		return -1
	}
	return i
}

// GetName method
func (h Host) GetName() string {
	return h.GetAttribute("NAME")
}

// GetImMad method
func (h Host) GetImMad() string {
	return h.GetAttribute("IM_MAD")
}

// GetVMMad method
func (h Host) GetVMMad() string {
	return h.GetAttribute("VM_MAD")
}

// GetCluster method
func (h Host) GetCluster() []Cluster {
	elements := h.XMLData.FindElements("CLUSTERS")
	clusters := make([]Cluster, len(elements))
	for i, e := range elements {
		clusters[i] = Cluster{XMLData: e}
	}
	return clusters
}
