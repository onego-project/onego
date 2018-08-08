package resources

import (
	"fmt"
	"github.com/beevik/etree"
	"strconv"
)

// DataStore struct
type DataStore struct {
	XMLData *etree.Element
}

// CreateDataStore constructs DataStore
func CreateDataStore(id int) *DataStore {
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0"`)

	el := doc.CreateElement("DATASTORE").CreateElement("ID")
	el.SetText(fmt.Sprintf("%d", id))

	return &DataStore{doc.Root()}
}

// GetAttribute method
func (d DataStore) GetAttribute(path string) string {
	elements := d.XMLData.FindElements(path)
	if elements == nil {
		return ""
	}
	return elements[0].Text()
}

// GetName method
func (d DataStore) GetName() string {
	return d.GetAttribute("NAME")
}

// GetClusters method
func (d DataStore) GetClusters() []Cluster {
	elements := d.XMLData.FindElements("CLUSTERS")
	clusters := make([]Cluster, len(elements))
	for i, e := range elements {
		clusters[i] = Cluster{XMLData: e}
	}
	return clusters
}

// GetID method
func (d DataStore) GetID() int {
	i, err := strconv.Atoi(d.GetAttribute("ID"))
	if err != nil {
		return -1
	}
	return i
}
