package requests

import "encoding/xml"

// ResizeRequest structure to create resize request for a virtual machine.
type ResizeRequest struct {
	XMLName xml.Name `xml:"TEMPLATE"`
	CPU     float64  `xml:"CPU"`
	VCpu    int      `xml:"VCPU"`
	Memory  int      `xml:"MEMORY"`
}
