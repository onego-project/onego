package services

import (
	"context"
	"encoding/xml"

	"github.com/onego-project/onego/resources"
)

// NetworkInterfaceService structure to attach/detach network interface to/from a virtual machine.
type NetworkInterfaceService struct {
	Service
}

type templateNIC struct {
	XMLName xml.Name       `xml:"TEMPLATE"`
	Nic     *resources.NIC `xml:"NIC"`
}

// Attach attaches a new network interface to the virtual machine.
func (nics *NetworkInterfaceService) Attach(ctx context.Context, vm resources.VirtualMachine, nic resources.NIC) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	nicText, err := resources.RenderInterfaceToXMLString(templateNIC{Nic: &nic})
	if err != nil {
		return err
	}

	_, err = nics.call(ctx, "one.vm.attachnic", vmID, nicText)

	return err
}

// Detach detaches a network interface from a virtual machine.
func (nics *NetworkInterfaceService) Detach(ctx context.Context, vm resources.VirtualMachine, nic resources.NIC) error {
	vmID, err := vm.ID()
	if err != nil {
		return err
	}

	_, err = nics.call(ctx, "one.vm.detachnic", vmID, nic.NicID)

	return err
}
