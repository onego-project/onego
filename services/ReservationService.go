package services

import (
	"context"

	"github.com/onego-project/onego/resources"
)

// ReservationService structure to manage reservation in OpenNebula virtual networks.
type ReservationService struct {
	Service
}

func (rs *ReservationService) reserve(ctx context.Context, parent resources.VirtualNetwork,
	reservation resources.Reservation) (*resources.VirtualNetwork, error) {
	vnID, err := parent.ID()
	if err != nil {
		return nil, err
	}

	reservationStr, err := resources.RenderInterfaceToXMLString(reservation)
	if err != nil {
		return nil, err
	}

	resArr, err := rs.call(ctx, "one.vn.reserve", vnID, reservationStr)
	if err != nil {
		return nil, err
	}

	vns := &VirtualNetworkService{Service: rs.Service}

	return vns.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Reserve reserves network addresses.
func (rs *ReservationService) Reserve(ctx context.Context, parent resources.VirtualNetwork,
	reservation resources.Reservation) (*resources.VirtualNetwork, error) {
	return rs.reserve(ctx, parent, reservation)
}

// ReserveToExisting reserves network addresses to existing reservation (virtual network).
func (rs *ReservationService) ReserveToExisting(ctx context.Context, parent resources.VirtualNetwork,
	size, addressRangeID, virtualNetworkID int) (*resources.VirtualNetwork, error) {
	return rs.reserve(ctx, parent, resources.Reservation{Size: &size, AddressRangeID: &addressRangeID,
		VirtualNetworkID: &virtualNetworkID})
}

// ReserveAsNew reserves network addresses as new reservation (virtual network).
func (rs *ReservationService) ReserveAsNew(ctx context.Context, parent resources.VirtualNetwork, size int,
	name string) (*resources.VirtualNetwork, error) {
	return rs.reserve(ctx, parent, resources.Reservation{Size: &size, Name: name})
}

// Free frees a reserved address range from a virtual network.
func (rs *ReservationService) Free(ctx context.Context, parent resources.VirtualNetwork,
	addressRange resources.AddressRange) error {
	vnID, err := parent.ID()
	if err != nil {
		return err
	}

	arid := addressRange.ID

	_, err = rs.call(ctx, "one.vn.free_ar", vnID, *arid)

	return err
}
