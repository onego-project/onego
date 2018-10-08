package services

import (
	"context"

	"github.com/beevik/etree"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/requests"
	"github.com/onego-project/onego/resources"
)

// ImageService structure to manage OpenNebula image.
type ImageService struct {
	Service
}

const pageOffsetDefault = -1
const pageSizeDefault = -1

// Allocate to create a new image in OpenNebula.
func (is *ImageService) Allocate(ctx context.Context, blueprint blueprint.Interface,
	datastore resources.Datastore) (*resources.Image, error) {
	datastoreID, err := datastore.ID()
	if err != nil {
		return nil, err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return nil, err
	}

	resArr, err := is.call(ctx, "one.image.allocate", blueprintText, datastoreID)
	if err != nil {
		return nil, err
	}

	return is.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Clone clones an existing image.
func (is *ImageService) Clone(ctx context.Context, image resources.Image, name string,
	datastore resources.Datastore) (*resources.Image, error) {
	imageID, err := image.ID()
	if err != nil {
		return nil, err
	}

	datastoreID, err := datastore.ID()
	if err != nil {
		return nil, err
	}

	resArr, err := is.call(ctx, "one.image.clone", imageID, name, datastoreID)
	if err != nil {
		return nil, err
	}

	return is.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Delete deletes the given image from the pool.
func (is *ImageService) Delete(ctx context.Context, image resources.Image) error {
	imageID, err := image.ID()
	if err != nil {
		return err
	}

	_, err = is.call(ctx, "one.image.delete", imageID)

	return err
}

func (is *ImageService) enable(ctx context.Context, image resources.Image, enable bool) error {
	imageID, err := image.ID()
	if err != nil {
		return err
	}

	_, err = is.call(ctx, "one.image.enable", imageID, enable)

	return err
}

// Enable enables an image.
func (is *ImageService) Enable(ctx context.Context, image resources.Image) error {
	return is.enable(ctx, image, true)
}

// Disable disables an image.
func (is *ImageService) Disable(ctx context.Context, image resources.Image) error {
	return is.enable(ctx, image, false)
}

func (is *ImageService) persistent(ctx context.Context, image resources.Image, persistent bool) error {
	imageID, err := image.ID()
	if err != nil {
		return err
	}

	_, err = is.call(ctx, "one.image.persistent", imageID, persistent)

	return err
}

// MakePersistent sets the Image as persistent.
func (is *ImageService) MakePersistent(ctx context.Context, image resources.Image) error {
	return is.persistent(ctx, image, true)
}

// MakeNonPersistent sets the Image as not persistent.
func (is *ImageService) MakeNonPersistent(ctx context.Context, image resources.Image) error {
	return is.persistent(ctx, image, false)
}

// ChangeType changes the type of an Image.
func (is *ImageService) ChangeType(ctx context.Context, image resources.Image, imageType resources.ImageType) error {
	imageID, err := image.ID()
	if err != nil {
		return err
	}

	_, err = is.call(ctx, "one.image.chtype", imageID, resources.ImageTypeMap[imageType])

	return err
}

// Chmod changes the permission bits of an image.
func (is *ImageService) Chmod(ctx context.Context, image resources.Image, request requests.PermissionRequest) error {
	imageID, err := image.ID()
	if err != nil {
		return err
	}

	_, err = is.call(ctx, "one.image.chmod", imageID,
		request.Permissions[requests.User][requests.Use], request.Permissions[requests.User][requests.Manage],
		request.Permissions[requests.User][requests.Admin],
		request.Permissions[requests.Group][requests.Use], request.Permissions[requests.Group][requests.Manage],
		request.Permissions[requests.Group][requests.Admin],
		request.Permissions[requests.Other][requests.Use], request.Permissions[requests.Other][requests.Manage],
		request.Permissions[requests.Other][requests.Admin])

	return err
}

// Chown changes the ownership of an image.
func (is *ImageService) Chown(ctx context.Context, image resources.Image, request requests.OwnershipRequest) error {
	imageID, err := image.ID()
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

	_, err = is.call(ctx, "one.image.chown", imageID, userID, groupID)

	return err
}

// Rename renames an image.
func (is *ImageService) Rename(ctx context.Context, image resources.Image, name string) error {
	imageID, err := image.ID()
	if err != nil {
		return err
	}

	_, err = is.call(ctx, "one.image.rename", imageID, name)

	return err
}

// Update replaces the image template contents.
func (is *ImageService) Update(ctx context.Context, image resources.Image, blueprint blueprint.Interface,
	updateType UpdateType) error {
	imageID, err := image.ID()
	if err != nil {
		return err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return err
	}

	_, err = is.call(ctx, "one.image.update", imageID, blueprintText, updateType)

	return err
}

// DeleteSnapshot deletes a snapshot from the image.
func (is *ImageService) DeleteSnapshot(ctx context.Context, image resources.Image,
	snapshot resources.ImageSnapshot) error {
	imageID, err := image.ID()
	if err != nil {
		return err
	}

	_, err = is.call(ctx, "one.image.snapshotdelete", imageID, snapshot.ID)

	return err
}

// RevertSnapshot reverts image state to a previous snapshot.
func (is *ImageService) RevertSnapshot(ctx context.Context, image resources.Image,
	snapshot resources.ImageSnapshot) error {
	imageID, err := image.ID()
	if err != nil {
		return err
	}

	_, err = is.call(ctx, "one.image.snapshotrevert", imageID, snapshot.ID)

	return err
}

// FlattenSnapshot flattens the snapshot of image and discards others.
func (is *ImageService) FlattenSnapshot(ctx context.Context, image resources.Image,
	snapshot resources.ImageSnapshot) error {
	imageID, err := image.ID()
	if err != nil {
		return err
	}

	_, err = is.call(ctx, "one.image.snapshotflatten", imageID, snapshot.ID)

	return err
}

// RetrieveInfo retrieves information for the image.
func (is *ImageService) RetrieveInfo(ctx context.Context, imageID int) (*resources.Image, error) {
	doc, err := is.retrieveInfo(ctx, "one.image.info", imageID)
	if err != nil {
		return nil, err
	}

	return resources.CreateImageFromXML(doc.Root()), nil
}

func (is *ImageService) list(ctx context.Context, filterFlag, pageOffset, pageSize int) ([]*resources.Image, error) {
	resArr, err := is.call(ctx, "one.imagepool.info", filterFlag, pageOffset, pageSize)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromString(resArr[resultIndex].ResultString()); err != nil {
		return nil, err
	}

	elements := doc.FindElements("IMAGE_POOL/IMAGE")

	images := make([]*resources.Image, len(elements))
	for i, e := range elements {
		images[i] = resources.CreateImageFromXML(e)
	}

	return images, nil
}

// ListAll retrieves information for part of the images in the pool which belong to given owner(s) in ownership filter.
func (is *ImageService) ListAll(ctx context.Context, filter OwnershipFilter) ([]*resources.Image, error) {
	return is.list(ctx, int(filter), pageOffsetDefault, pageSizeDefault)
}

// ListAllForUser retrieves information for part of the images in the pool.
func (is *ImageService) ListAllForUser(ctx context.Context, user resources.User) ([]*resources.Image, error) {
	userID, err := user.ID()
	if err != nil {
		return nil, err
	}

	return is.list(ctx, userID, pageOffsetDefault, pageSizeDefault)
}

// List retrieves information for all the images in the pool.
func (is *ImageService) List(ctx context.Context, pageOffset int, pageSize int,
	filter OwnershipFilter) ([]*resources.Image, error) {
	return is.list(ctx, int(filter), (pageOffset-1)*pageSize, -pageSize)
}

// ListForUser retrieves information for part of the images in the pool.
func (is *ImageService) ListForUser(ctx context.Context, user resources.User, pageOffset int,
	pageSize int) ([]*resources.Image, error) {
	userID, err := user.ID()
	if err != nil {
		return nil, err
	}

	return is.list(ctx, userID, (pageOffset-1)*pageSize, -pageSize)
}
