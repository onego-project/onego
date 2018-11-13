package services

import (
	"context"
	"encoding/xml"
	"net"

	"github.com/onego-project/onego/resources"
)

// LeaseService structure to manage lease in OpenNebula virtual networks.
type LeaseService struct {
	Service
}

type templateLease struct {
	XMLName    xml.Name           `xml:"TEMPLATE,omitempty"`
	LeaseMngmt []*leaseManagement `xml:"LEASES,omitempty"`
}

type leaseManagement struct {
	XMLName xml.Name `xml:"LEASES,omitempty"`
	IP      *net.IP  `xml:"IP,omitempty"`
}

func (ls *LeaseService) manageLease(ctx context.Context, methodName string, vn resources.VirtualNetwork,
	ips []net.IP) error {
	vnID, err := vn.ID()
	if err != nil {
		return err
	}

	leases := make([]*leaseManagement, len(ips))
	for i, ip := range ips {
		leases[i] = &leaseManagement{IP: &ip}
	}

	leaseText, err := resources.RenderInterfaceToXMLString(templateLease{LeaseMngmt: leases})
	if err != nil {
		return err
	}

	_, err = ls.call(ctx, methodName, vnID, leaseText)

	return err
}

// Hold holds a virtual network Lease as used.
func (ls *LeaseService) Hold(ctx context.Context, vn resources.VirtualNetwork, ips []net.IP) error {
	return ls.manageLease(ctx, "one.vn.hold", vn, ips)
}

// Release releases a virtual network Lease on hold.
func (ls *LeaseService) Release(ctx context.Context, vn resources.VirtualNetwork, ips []net.IP) error {
	return ls.manageLease(ctx, "one.vn.release", vn, ips)
}
