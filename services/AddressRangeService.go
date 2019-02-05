package services

import (
	"context"
	"encoding/xml"

	"github.com/onego-project/onego/errors"

	"github.com/onego-project/onego/resources"
)

// AddressRangeService structure to manage address range in OpenNebula virtual networks.
type AddressRangeService struct {
	Service
}

type templateAddressRange struct {
	XMLName xml.Name               `xml:"TEMPLATE"`
	AR      resources.AddressRange `xml:"AR"`
}

func (ars *AddressRangeService) manageAddressRange(ctx context.Context, methodName string, vn resources.VirtualNetwork,
	ar resources.AddressRange) (*resources.AddressRange, error) {
	vnID, err := vn.ID()
	if err != nil {
		return nil, err
	}

	arText, err := resources.RenderInterfaceToXMLString(templateAddressRange{AR: ar})
	if err != nil {
		// error should never occur
		return nil, err
	}

	resArr, err := ars.call(ctx, methodName, vnID, arText)
	if err != nil {
		return nil, err
	}

	vns := &VirtualNetworkService{Service: ars.Service}

	vnet, err := vns.RetrieveInfo(ctx, vnID)
	if err != nil {
		// error should never occur (only in case when vnet is deleted during run this method)
		return nil, err
	}

	arArray, err := vnet.AddressRanges()
	if err != nil {
		// error should never occur (only if vnet was modified from outside and doesn't contain some of attributes)
		return nil, err
	}

	for _, e := range arArray {
		id := e.ID
		if *id == int(resArr[1].ResultInt()) {
			return e, nil
		}
	}

	// error should never occur
	return nil, errors.ErrAddressRangeSetWrong
}

// Add adds address range to virtual network.
// AR must contain TYPE, SIZE, IP.
func (ars *AddressRangeService) Add(ctx context.Context, vn resources.VirtualNetwork,
	ar resources.AddressRange) (*resources.AddressRange, error) {
	return ars.manageAddressRange(ctx, "one.vn.add_ar", vn, ar)
}

// Update updates the attributes of an address range.
func (ars *AddressRangeService) Update(ctx context.Context, vn resources.VirtualNetwork,
	ar resources.AddressRange) (*resources.AddressRange, error) {
	return ars.manageAddressRange(ctx, "one.vn.update_ar", vn, ar)
}

// Delete removes an address range from a virtual network.
func (ars *AddressRangeService) Delete(ctx context.Context, vn resources.VirtualNetwork,
	ar resources.AddressRange) error {
	vnID, err := vn.ID()
	if err != nil {
		return err
	}

	if ar.ID == nil {
		return errors.ErrAddressRangeNoID
	}

	_, err = ars.call(ctx, "one.vn.rm_ar", vnID, *ar.ID)

	return err
}
