package services

import (
	"context"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/requests"
	"github.com/onego-project/onego/resources"
)

// DatastoreService structure to manage OpenNebula datastore
type DatastoreService struct {
	Service
}

// Allocate to create a new datastore in OpenNebula
func (ds *DatastoreService) Allocate(ctx context.Context, blueprint blueprint.Interface, cluster resources.Cluster) (*resources.Datastore, error) {
	clusterID, err := cluster.ID()
	if err != nil {
		return nil, err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return nil, err
	}

	resArr, err := ds.call(ctx, "one.datastore.allocate", blueprintText, clusterID)
	if err != nil {
		return nil, err
	}

	return ds.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Delete to delete the given datastore from the pool
func (ds *DatastoreService) Delete(ctx context.Context, datastore resources.Datastore) error {
	datastoreID, err := datastore.ID()
	if err != nil {
		return err
	}

	_, err = ds.call(ctx, "one.datastore.delete", datastoreID)

	return err
}

// Update to replace the datastore template contents
// Update types: Replace or Merge
func (ds *DatastoreService) Update(ctx context.Context, datastore resources.Datastore, blueprint blueprint.Interface, updateType UpdateType) error {
	datastoreID, err := datastore.ID()
	if err != nil {
		return err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return err
	}

	_, err = ds.call(ctx, "one.datastore.update", datastoreID, blueprintText, updateType)

	return err
}

// Chmod to change the permission bits of a datastore
// If bit is set to -1, it will not change.
func (ds *DatastoreService) Chmod(ctx context.Context, datastore resources.Datastore, request requests.PermissionRequest) error {
	datastoreID, err := datastore.ID()
	if err != nil {
		return err
	}

	_, err = ds.call(ctx, "one.datastore.chmod", datastoreID,
		request.Permissions[requests.User][requests.Use], request.Permissions[requests.User][requests.Manage], request.Permissions[requests.User][requests.Admin],
		request.Permissions[requests.Group][requests.Use], request.Permissions[requests.Group][requests.Manage], request.Permissions[requests.Group][requests.Admin],
		request.Permissions[requests.Other][requests.Use], request.Permissions[requests.Other][requests.Manage], request.Permissions[requests.Other][requests.Admin])

	return err
}

// Chown to change the ownership of the datastore
// If user id is set to -1, the owner is not changed.
// If group id is set to -1, the group is not changed.
func (ds *DatastoreService) Chown(ctx context.Context, datastore resources.Datastore, request requests.OwnershipRequest) error {
	datastoreID, err := datastore.ID()
	if err != nil {
		return err
	}

	userID, err := request.User.ID()
	if err != nil {
		return err
	}

	groupID, err := request.Group.ID()
	if err != nil {
		return err
	}

	_, err = ds.call(ctx, "one.datastore.chown", datastoreID, userID, groupID)

	return err
}

// Rename to rename a datastore
func (ds *DatastoreService) Rename(ctx context.Context, datastore resources.Datastore, name string) error {
	datastoreID, err := datastore.ID()
	if err != nil {
		return err
	}

	_, err = ds.call(ctx, "one.datastore.rename", datastoreID, name)

	return err
}

func (ds *DatastoreService) enable(ctx context.Context, datastore resources.Datastore, enable bool) error {
	datastoreID, err := datastore.ID()
	if err != nil {
		return err
	}

	_, err = ds.call(ctx, "one.datastore.enable", datastoreID, enable)

	return err
}

// Enable to enable a datastore
func (ds *DatastoreService) Enable(ctx context.Context, datastore resources.Datastore) error {
	return ds.enable(ctx, datastore, true)
}

// Disable to disable a datastore
func (ds *DatastoreService) Disable(ctx context.Context, datastore resources.Datastore) error {
	return ds.enable(ctx, datastore, false)
}

// RetrieveInfo to retrieve information for the datastore
func (ds *DatastoreService) RetrieveInfo(ctx context.Context, datastoreID int) (*resources.Datastore, error) {
	doc, err := ds.retrieveInfo(ctx, "one.datastore.info", datastoreID)
	if err != nil {
		return nil, err
	}

	return resources.CreateDatastoreFromXML(doc.Root()), nil
}

// List to retrieve information for all of the datastores in the pool
func (ds *DatastoreService) List(ctx context.Context) ([]*resources.Datastore, error) {
	doc, err := ds.list(ctx, "one.datastorepool.info")
	if err != nil {
		return nil, err
	}

	elements := doc.FindElements("DATASTORE_POOL/DATASTORE")

	datastores := make([]*resources.Datastore, len(elements))
	for i, e := range elements {
		datastores[i] = resources.CreateDatastoreFromXML(e)
	}

	return datastores, nil
}
