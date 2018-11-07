package services

import (
	"context"

	"github.com/beevik/etree"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/requests"
	"github.com/onego-project/onego/resources"
)

// VirtualNetworkService structure to manage OpenNebula virtual network.
type VirtualNetworkService struct {
	Service
}

// Allocate creates a new virtual network in OpenNebula.
func (vns *VirtualNetworkService) Allocate(ctx context.Context, vnBlueprint blueprint.Interface,
	cluster resources.Cluster) (*resources.VirtualNetwork, error) {
	clusterID, err := cluster.ID()
	if err != nil {
		return nil, err
	}

	blueprintText, err := vnBlueprint.Render()
	if err != nil {
		return nil, err
	}

	resArr, err := vns.call(ctx, "one.vn.allocate", blueprintText, clusterID)
	if err != nil {
		return nil, err
	}

	return vns.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Delete deletes the given virtual network from the pool.
func (vns *VirtualNetworkService) Delete(ctx context.Context, vn resources.VirtualNetwork) error {
	vnID, err := vn.ID()
	if err != nil {
		return err
	}

	_, err = vns.call(ctx, "one.vn.delete", vnID)

	return err
}

// Update replaces the virtual network template contents.
func (vns *VirtualNetworkService) Update(ctx context.Context, vn resources.VirtualNetwork,
	blueprint blueprint.Interface, updateType UpdateType) error {
	vnID, err := vn.ID()
	if err != nil {
		return err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return err
	}

	_, err = vns.call(ctx, "one.vn.update", vnID, blueprintText, updateType)

	return err
}

// Chmod to change the permission bits of a virtual network.
// If bit is set to -1, it will not change.
func (vns *VirtualNetworkService) Chmod(ctx context.Context, vn resources.VirtualNetwork,
	request requests.PermissionRequest) error {
	vnID, err := vn.ID()
	if err != nil {
		return err
	}

	return vns.chmod(ctx, "one.vn.chmod", vnID, request)
}

// Chown to change the ownership of the virtual network.
// If user id is set to -1, the owner is not changed.
// If group id is set to -1, the group is not changed.
func (vns *VirtualNetworkService) Chown(ctx context.Context, vn resources.VirtualNetwork,
	request requests.OwnershipRequest) error {
	vnID, err := vn.ID()
	if err != nil {
		return err
	}

	return vns.chown(ctx, "one.vn.chown", vnID, request)
}

// Rename renames a Virtual Network.
func (vns *VirtualNetworkService) Rename(ctx context.Context, vn resources.VirtualNetwork, name string) error {
	vnID, err := vn.ID()
	if err != nil {
		return err
	}

	_, err = vns.call(ctx, "one.vn.rename", vnID, name)

	return err
}

// RetrieveInfo retrieves information for the virtual network.
func (vns *VirtualNetworkService) RetrieveInfo(ctx context.Context, vnID int) (*resources.VirtualNetwork, error) {
	doc, err := vns.retrieveInfo(ctx, "one.vn.info", vnID)
	if err != nil {
		return nil, err
	}

	return resources.CreateVirtualNetworkFromXML(doc.Root()), nil
}

func (vns *VirtualNetworkService) list(ctx context.Context, filterFlag, pageOffset,
	pageSize int) ([]*resources.VirtualNetwork, error) {
	resArr, err := vns.call(ctx, "one.vnpool.info", filterFlag, pageOffset, pageSize)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromString(resArr[resultIndex].ResultString()); err != nil {
		return nil, err
	}

	elements := doc.FindElements("VNET_POOL/VNET")

	vnets := make([]*resources.VirtualNetwork, len(elements))
	for i, e := range elements {
		vnets[i] = resources.CreateVirtualNetworkFromXML(e)
	}

	return vnets, nil
}

// ListAll retrieves information for part of the virtual networks in the pool which belong to given owner(s)
// in ownership filter.
func (vns *VirtualNetworkService) ListAll(ctx context.Context,
	filter OwnershipFilter) ([]*resources.VirtualNetwork, error) {
	return vns.list(ctx, int(filter), pageOffsetDefault, pageSizeDefault)
}

// ListAllForUser retrieves information for part of the images in the pool.
func (vns *VirtualNetworkService) ListAllForUser(ctx context.Context,
	user resources.User) ([]*resources.VirtualNetwork, error) {
	userID, err := user.ID()
	if err != nil {
		return nil, err
	}

	return vns.list(ctx, userID, pageOffsetDefault, pageSizeDefault)
}

// List retrieves information for all the virtual networks in the pool.
func (vns *VirtualNetworkService) List(ctx context.Context, pageOffset int, pageSize int,
	filter OwnershipFilter) ([]*resources.VirtualNetwork, error) {
	return vns.list(ctx, int(filter), (pageOffset-1)*pageSize, -pageSize)
}

// ListForUser retrieves information for part of the virtual networks in the pool.
func (vns *VirtualNetworkService) ListForUser(ctx context.Context, user resources.User, pageOffset int,
	pageSize int) ([]*resources.VirtualNetwork, error) {
	userID, err := user.ID()
	if err != nil {
		return nil, err
	}

	return vns.list(ctx, userID, (pageOffset-1)*pageSize, -pageSize)
}
