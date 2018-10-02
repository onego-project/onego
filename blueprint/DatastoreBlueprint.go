package blueprint

import "github.com/onego-project/onego/resources"

// DatastoreBlueprint to set Datastore elements
type DatastoreBlueprint struct {
	Blueprint
}

// CreateAllocateDatastoreBlueprint creates empty DatastoreBlueprint
func CreateAllocateDatastoreBlueprint() *DatastoreBlueprint {
	return &DatastoreBlueprint{Blueprint: *CreateBlueprint("DATASTORE")}
}

// CreateUpdateDatastoreBlueprint creates empty DatastoreBlueprint
func CreateUpdateDatastoreBlueprint() *DatastoreBlueprint {
	return &DatastoreBlueprint{Blueprint: *CreateBlueprint("TEMPLATE")}
}

// SetDiskType sets disk type of the given datastore
func (ds *DatastoreBlueprint) SetDiskType(diskType resources.DiskType) {
	ds.SetElement("DISK_TYPE", resources.DiskTypeMap[diskType])
}

// SetDsMad sets DS_MAD of the given datastore
func (ds *DatastoreBlueprint) SetDsMad(dsMad string) {
	ds.SetElement("DS_MAD", dsMad)
}

// SetTmMad sets TM_MAD of the given datastore
func (ds *DatastoreBlueprint) SetTmMad(tmMad string) {
	ds.SetElement("TM_MAD", tmMad)
}

// SetType sets type of the given datastore
func (ds *DatastoreBlueprint) SetType(dsType resources.DatastoreType) {
	ds.SetElement("TYPE", resources.DatastoreTypeMap[dsType])
}
