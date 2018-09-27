package resources

import "github.com/beevik/etree"

// Image structure represents OpenNebula Image
type Image struct {
	Resource
}

// CreateImageWithID constructs Image with id
func CreateImageWithID(id int) *Image {
	return &Image{*CreateResource("IMAGE", id)}
}

// CreateImageFromXML constructs Image with full xml data
func CreateImageFromXML(XMLdata *etree.Element) *Image {
	return &Image{Resource: Resource{XMLData: XMLdata}}
}
