package errors

import (
	"errors"
	"fmt"
)

// OpenNebulaError structure represents errors caused by OpenNebula
type OpenNebulaError struct {
	Code     int
	Message  string
	ObjectID int
}

// XMLElementError structure represents errors in XML elements
type XMLElementError struct {
	Path string
}

// ErrNoClient error
var ErrNoClient = errors.New("no client")

// ErrNoImage error
var ErrNoImage = errors.New("no image to finish test")

// ErrNoImageBlueprint error
var ErrNoImageBlueprint = errors.New("no image blueprint to finish test")

// ErrNoGroup error
var ErrNoGroup = errors.New("no group to finish test")

// ErrNoDatastore error
var ErrNoDatastore = errors.New("no datastore to finish test")

// ErrNoDatastoreBlueprint error
var ErrNoDatastoreBlueprint = errors.New("no datastore blueprint to finish test")

// ErrNoUser error
var ErrNoUser = errors.New("no user to finish test")

// ErrNoUserBlueprint error
var ErrNoUserBlueprint = errors.New("no user blueprint to finish test")

// ErrBlueprintXMLEmpty error
var ErrBlueprintXMLEmpty = errors.New("blueprint XML data is empty")

// ErrNoHost error
var ErrNoHost = errors.New("no host to finish test")

// ErrNoHostBlueprint error
var ErrNoHostBlueprint = errors.New("no host blueprint to finish test")

// NoObjectID to distinguish errors from OpenNebula with 3 or 4 arguments
var NoObjectID = -1

func (one *OpenNebulaError) Error() string {
	if one.ObjectID != NoObjectID {
		return fmt.Sprintf("%s, error code: %d, id of the object "+
			"that caused the error %d", one.Message, one.Code, one.ObjectID)
	}
	return fmt.Sprintf("%s, error code: %d", one.Message, one.Code)
}

func (xee *XMLElementError) Error() string {
	return fmt.Sprintf("no element %s", xee.Path)
}