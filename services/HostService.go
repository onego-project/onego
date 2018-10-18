package services

import (
	"context"

	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/resources"
)

// HostService structure to manage OpenNebula host.
type HostService struct {
	Service
}

// HostStatus to set Host status
type HostStatus int

const (
	// HostEnabled to enable Host
	HostEnabled HostStatus = iota
	// HostDisabled to disable Host
	HostDisabled
	// HostOffline to power off Host
	HostOffline
)

// Allocate creates a new host in OpenNebula.
func (hs *HostService) Allocate(ctx context.Context, hostName, imMad, vmMad string,
	cluster resources.Cluster) (*resources.Host, error) {
	clusterID, err := cluster.ID()
	if err != nil {
		return nil, err
	}

	resArr, err := hs.call(ctx, "one.host.allocate", hostName, imMad, vmMad, clusterID)

	if err != nil {
		return nil, err
	}

	return hs.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Delete deletes the given host from the pool.
func (hs *HostService) Delete(ctx context.Context, host resources.Host) error {
	hostID, err := host.ID()
	if err != nil {
		return err
	}

	_, err = hs.call(ctx, "one.host.delete", hostID)

	return err
}

func (hs *HostService) status(ctx context.Context, host resources.Host, hostStatus HostStatus) error {
	hostID, err := host.ID()
	if err != nil {
		return err
	}

	_, err = hs.call(ctx, "one.host.status", hostID, hostStatus)

	return err
}

// Enable enables host.
func (hs *HostService) Enable(ctx context.Context, host resources.Host) error {
	return hs.status(ctx, host, HostEnabled)
}

// Disable disables host.
func (hs *HostService) Disable(ctx context.Context, host resources.Host) error {
	return hs.status(ctx, host, HostDisabled)
}

// Offline power off host.
func (hs *HostService) Offline(ctx context.Context, host resources.Host) error {
	return hs.status(ctx, host, HostOffline)
}

// Update replaces the host template contents.
func (hs *HostService) Update(ctx context.Context, host resources.Host, blueprint blueprint.Interface,
	updateType UpdateType) error {
	hostID, err := host.ID()
	if err != nil {
		return err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return err
	}

	_, err = hs.call(ctx, "one.host.update", hostID, blueprintText, updateType)

	return err
}

// Rename renames given host.
func (hs *HostService) Rename(ctx context.Context, host resources.Host, name string) error {
	hostID, err := host.ID()
	if err != nil {
		return err
	}

	_, err = hs.call(ctx, "one.host.rename", hostID, name)

	return err
}

// RetrieveInfo retrieves information for the host.
func (hs *HostService) RetrieveInfo(ctx context.Context, hostID int) (*resources.Host, error) {
	doc, err := hs.retrieveInfo(ctx, "one.host.info", hostID)
	if err != nil {
		return nil, err
	}

	return resources.CreateHostFromXML(doc.Root()), nil
}

// List to retrieve information for all of the host in the pool
func (hs *HostService) List(ctx context.Context) ([]*resources.Host, error) {
	doc, err := hs.list(ctx, "one.hostpool.info")
	if err != nil {
		return nil, err
	}

	elements := doc.FindElements("HOST_POOL/HOST")

	hosts := make([]*resources.Host, len(elements))
	for i, e := range elements {
		hosts[i] = resources.CreateHostFromXML(e)
	}

	return hosts, nil
}
