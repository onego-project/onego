package resources

import (
	"encoding/xml"
	"net"

	"github.com/beevik/etree"
	"github.com/onego-project/onego/errors"
)

// VirtualNetwork structure to manage OpenNebula Virtual Network
type VirtualNetwork struct {
	Resource
}

// AddressRange structure represents Address Range in Virtual Network
type AddressRange struct {
	XMLName    xml.Name `xml:"AR,omitempty"`
	ID         *int     `xml:"AR_ID,omitempty"`
	IP         *net.IP  `xml:"IP,omitempty"`
	Mac        string   `xml:"MAC,omitempty"`
	Size       *int     `xml:"SIZE,omitempty"`
	Type       string   `xml:"TYPE,omitempty"`
	MacEnd     string   `xml:"MAC_END,omitempty"`
	IPEnd      net.IP   `xml:"IP_END,omitempty"`
	UsedLeases *int     `xml:"USED_LEASES,omitempty"`
	Leases     []*Lease `xml:"LEASES,omitempty"`
}

// Reservation structure to reserve network address in OpenNebula virtual network.
type Reservation struct {
	XMLName          xml.Name `xml:"TEMPLATE,omitempty"`
	Size             *int     `xml:"SIZE,omitempty"`
	Name             string   `xml:"NAME,omitempty"`
	AddressRangeID   *int     `xml:"AR_ID,omitempty"`
	VirtualNetworkID *int     `xml:"NETWORK_ID,omitempty"`
	Mac              string   `xml:"MAC,omitempty"`
	IP               *net.IP  `xml:"IP,omitempty"`
}

// Lease structure represents Lease in Address Range
type Lease struct {
	XMLName          xml.Name `xml:"LEASE,omitempty"`
	IP               *net.IP  `xml:"IP,omitempty"`
	Mac              string   `xml:"MAC,omitempty"`
	VirtualMachineID *int     `xml:"VM,omitempty"`
}

// CreateVirtualNetworkWithID constructs VirtualNetwork with id
func CreateVirtualNetworkWithID(id int) *VirtualNetwork {
	return &VirtualNetwork{*CreateResource("VNET", id)}
}

// CreateVirtualNetworkFromXML constructs VirtualNetwork with full xml data
func CreateVirtualNetworkFromXML(XMLdata *etree.Element) *VirtualNetwork {
	return &VirtualNetwork{Resource: Resource{XMLData: XMLdata}}
}

// User gets user ID of given virtual network
func (vn *VirtualNetwork) User() (int, error) {
	return vn.intAttribute("UID")
}

// Group gets group ID of given virtual network
func (vn *VirtualNetwork) Group() (int, error) {
	return vn.intAttribute("GID")
}

// Permissions gets Permissions structure with permissions of given virtual network
func (vn *VirtualNetwork) Permissions() (*Permissions, error) {
	return vn.permissions()
}

// Clusters gets array of cluster IDs of given virtual network
func (vn *VirtualNetwork) Clusters() ([]int, error) {
	return vn.arrayOfIDs("CLUSTERS")
}

// Bridge gets bridge of given virtual network
func (vn *VirtualNetwork) Bridge() (string, error) {
	return vn.Attribute("BRIDGE")
}

// ParentNetworkID gets id of parent network
func (vn *VirtualNetwork) ParentNetworkID() (int, error) {
	return vn.intAttribute("PARENT_NETWORK_ID")
}

// VnMad gets VN_MAD of given virtual network
func (vn *VirtualNetwork) VnMad() (string, error) {
	return vn.Attribute("VN_MAD")
}

// PhysicalDevice gets physical device of given virtual network
func (vn *VirtualNetwork) PhysicalDevice() (string, error) {
	return vn.Attribute("PHYDEV")
}

// VirtualLanID gets virtual lan ID of given virtual network
func (vn *VirtualNetwork) VirtualLanID() (int, error) {
	return vn.intAttribute("VLAN_ID")
}

// VirtualLanIDAutomatic gets automatic virtual lan ID of given virtual network
func (vn *VirtualNetwork) VirtualLanIDAutomatic() (int, error) {
	return vn.intAttribute("VLAN_ID_AUTOMATIC")
}

// UsedLeases gets number of used leases of given virtual network
func (vn *VirtualNetwork) UsedLeases() (int, error) {
	return vn.intAttribute("USED_LEASES")
}

// VirtualRouters gets array of virtual routers IDs of given virtual network
func (vn *VirtualNetwork) VirtualRouters() ([]int, error) {
	return vn.arrayOfIDs("VROUTERS")
}

// AddressRanges gets array of Address ranges of given virtual network
func (vn *VirtualNetwork) AddressRanges() ([]*AddressRange, error) {
	elements := vn.XMLData.FindElements("AR_POOL/AR")
	if len(elements) == 0 {
		return make([]*AddressRange, 0), nil
	}

	ars := make([]*AddressRange, len(elements))
	var err error

	for i, e := range elements {
		ars[i], err = createAddressRangeFromElement(e)
		if err != nil {
			return nil, err
		}
	}
	return ars, nil
}

func createAddressRangeFromElement(element *etree.Element) (*AddressRange, error) {
	if element == nil {
		return nil, &errors.XMLElementError{Path: "AR_POOL/AR"}
	}

	parseStrings := []string{"MAC", "TYPE", "MAC_END", "IP", "IP_END"}
	parsedStrings, err := parseStringsFromElement(element, parseStrings)
	if err != nil {
		return nil, err
	}

	parseInts := []string{"AR_ID", "SIZE", "USED_LEASES"}
	parsedInts, err := parseIntsFromElement(element, parseInts)
	if err != nil {
		return nil, err
	}

	parsedLeases, err := leases(element)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(parsedStrings[3])

	return &AddressRange{ID: &parsedInts[0], IP: &ip, Mac: parsedStrings[0],
		Size: &parsedInts[1], Type: parsedStrings[1], MacEnd: parsedStrings[2], IPEnd: net.ParseIP(parsedStrings[4]),
		UsedLeases: &parsedInts[2], Leases: parsedLeases}, nil
}

func leases(element *etree.Element) ([]*Lease, error) {
	elements := element.FindElements("LEASES/LEASE")
	if len(elements) == 0 {
		return make([]*Lease, 0), nil
	}

	leases := make([]*Lease, len(elements))
	var err error

	for i, e := range elements {
		leases[i], err = createLeaseFromElement(e)
		if err != nil {
			return nil, err
		}
	}
	return leases, nil
}

func createLeaseFromElement(element *etree.Element) (*Lease, error) {
	if element == nil {
		return nil, &errors.XMLElementError{Path: "AR_POOL/AR/LEASES/LEASE"}
	}

	ipString, err := attributeFromElement(element, "IP")
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(ipString)

	mac, err := attributeFromElement(element, "MAC")
	if err != nil {
		return nil, err
	}

	vm, err := intAttributeFromElement(element, "VM")
	if err != nil {
		// not necessary attribute
		return &Lease{IP: &ip, Mac: mac}, nil
	}

	return &Lease{IP: &ip, Mac: mac, VirtualMachineID: &vm}, nil
}
