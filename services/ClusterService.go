package services

import (
	"context"

	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/resources"
)

// ClusterService struct
type ClusterService struct {
	Service
}

// Allocate creates a new cluster in OpenNebula.
func (cs *ClusterService) Allocate(ctx context.Context, clusterName string) (*resources.Cluster, error) {
	resArr, err := cs.call(ctx, "one.cluster.allocate", clusterName)

	if err != nil {
		return nil, err
	}

	return cs.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Delete deletes the given cluster from the pool.
func (cs *ClusterService) Delete(ctx context.Context, cluster resources.Cluster) error {
	clusterID, err := cluster.ID()
	if err != nil {
		return err
	}

	_, err = cs.call(ctx, "one.cluster.delete", clusterID)

	return err
}

// Update replaces the cluster template contents.
func (cs *ClusterService) Update(ctx context.Context, cluster resources.Cluster, blueprint blueprint.Interface,
	updateType UpdateType) error {
	clusterID, err := cluster.ID()
	if err != nil {
		return err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return err
	}

	_, err = cs.call(ctx, "one.cluster.update", clusterID, blueprintText, updateType)

	return err
}

// AddHost adds a host to the given cluster
func (cs *ClusterService) AddHost(ctx context.Context, cluster resources.Cluster, host resources.Host) error {
	clusterID, err := cluster.ID()
	if err != nil {
		return err
	}

	hostID, err := host.ID()
	if err != nil {
		return err
	}

	_, err = cs.call(ctx, "one.cluster.addhost", clusterID, hostID)

	return err
}

// DeleteHost deletes a host to the given cluster
func (cs *ClusterService) DeleteHost(ctx context.Context, cluster resources.Cluster, host resources.Host) error {
	clusterID, err := cluster.ID()
	if err != nil {
		return err
	}

	hostID, err := host.ID()
	if err != nil {
		return err
	}

	_, err = cs.call(ctx, "one.cluster.delhost", clusterID, hostID)

	return err
}

// AddDatastore adds a datastore to the given cluster
func (cs *ClusterService) AddDatastore(ctx context.Context, cluster resources.Cluster,
	datastore resources.Datastore) error {
	clusterID, err := cluster.ID()
	if err != nil {
		return err
	}

	datastoreID, err := datastore.ID()
	if err != nil {
		return err
	}

	_, err = cs.call(ctx, "one.cluster.adddatastore", clusterID, datastoreID)

	return err
}

// DeleteDatastore deletes a datastore to the given cluster
func (cs *ClusterService) DeleteDatastore(ctx context.Context, cluster resources.Cluster,
	datastore resources.Datastore) error {
	clusterID, err := cluster.ID()
	if err != nil {
		return err
	}

	datastoreID, err := datastore.ID()
	if err != nil {
		return err
	}

	_, err = cs.call(ctx, "one.cluster.deldatastore", clusterID, datastoreID)

	return err
}

// AddVirtualNetwork adds a virtual network to the given cluster
func (cs *ClusterService) AddVirtualNetwork(ctx context.Context, cluster resources.Cluster,
	virtualNetwork resources.VirtualNetwork) error {
	clusterID, err := cluster.ID()
	if err != nil {
		return err
	}

	virtualNetworkID, err := virtualNetwork.ID()
	if err != nil {
		return err
	}

	_, err = cs.call(ctx, "one.cluster.addvnet", clusterID, virtualNetworkID)

	return err
}

// DeleteVirtualNetwork deletes a virtual network to the given cluster
func (cs *ClusterService) DeleteVirtualNetwork(ctx context.Context, cluster resources.Cluster,
	virtualNetwork resources.VirtualNetwork) error {
	clusterID, err := cluster.ID()
	if err != nil {
		return err
	}

	virtualNetworkID, err := virtualNetwork.ID()
	if err != nil {
		return err
	}

	_, err = cs.call(ctx, "one.cluster.delvnet", clusterID, virtualNetworkID)

	return err
}

// Rename renames given cluster.
func (cs *ClusterService) Rename(ctx context.Context, cluster resources.Cluster, name string) error {
	clusterID, err := cluster.ID()
	if err != nil {
		return err
	}

	_, err = cs.call(ctx, "one.cluster.rename", clusterID, name)

	return err
}

// RetrieveInfo retrieves information for the cluster.
func (cs *ClusterService) RetrieveInfo(ctx context.Context, clusterID int) (*resources.Cluster, error) {
	doc, err := cs.retrieveInfo(ctx, "one.cluster.info", clusterID)
	if err != nil {
		return nil, err
	}

	return resources.CreateClusterFromXML(doc.Root()), nil
}

// List to retrieve information for all of the cluster in the pool
func (cs *ClusterService) List(ctx context.Context) ([]*resources.Cluster, error) {
	doc, err := cs.list(ctx, "one.clusterpool.info")
	if err != nil {
		return nil, err
	}

	elements := doc.FindElements("CLUSTER_POOL/CLUSTER")

	clusters := make([]*resources.Cluster, len(elements))
	for i, e := range elements {
		clusters[i] = resources.CreateClusterFromXML(e)
	}

	return clusters, nil
}
