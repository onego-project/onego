package services

import (
	"context"

	"github.com/beevik/etree"

	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/requests"
	"github.com/onego-project/onego/resources"
)

// VMTemplateService structure to manage OpenNebula Template.
type VMTemplateService struct {
	Service
}

// Allocate allocates a new template in OpenNebula.
func (vmts *VMTemplateService) Allocate(ctx context.Context, p blueprint.Interface) (*resources.VMTemplate, error) {
	blueprintText, err := p.Render()
	if err != nil {
		return nil, err
	}

	resArr, err := vmts.call(ctx, "one.template.allocate", blueprintText)
	if err != nil {
		return nil, err
	}

	return vmts.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Clone clones an existing virtual machine template.
func (vmts *VMTemplateService) Clone(ctx context.Context, template resources.VMTemplate, name string,
	withImages bool) (*resources.VMTemplate, error) {
	templateID, err := template.ID()
	if err != nil {
		return nil, err
	}

	resArr, err := vmts.call(ctx, "one.template.clone", templateID, name, withImages)
	if err != nil {
		return nil, err
	}

	return vmts.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Delete deletes the given template from the pool.
func (vmts *VMTemplateService) Delete(ctx context.Context, template resources.VMTemplate, withImages bool) error {
	templateID, err := template.ID()
	if err != nil {
		return err
	}

	_, err = vmts.call(ctx, "one.template.delete", templateID, withImages)

	return err
}

// Update merges or replaces the template contents.
func (vmts *VMTemplateService) Update(ctx context.Context, template resources.VMTemplate, p blueprint.Interface,
	updateType UpdateType) error {
	templateID, err := template.ID()
	if err != nil {
		return err
	}

	blueprintText, err := p.Render()
	if err != nil {
		return err
	}

	_, err = vmts.call(ctx, "one.template.update", templateID, blueprintText, updateType)

	return err
}

// Chmod changes the permission bits of a template.
func (vmts *VMTemplateService) Chmod(ctx context.Context, template resources.VMTemplate,
	request requests.PermissionRequest) error {
	templateID, err := template.ID()
	if err != nil {
		return err
	}

	return vmts.chmod(ctx, "one.template.chmod", templateID, request)
}

// Chown changes the ownership of a template.
func (vmts *VMTemplateService) Chown(ctx context.Context, template resources.VMTemplate,
	request requests.OwnershipRequest) error {
	templateID, err := template.ID()
	if err != nil {
		return err
	}

	return vmts.chown(ctx, "one.template.chown", templateID, request)
}

// Rename renames a template.
func (vmts *VMTemplateService) Rename(ctx context.Context, template resources.VMTemplate, name string) error {
	templateID, err := template.ID()
	if err != nil {
		return err
	}

	_, err = vmts.call(ctx, "one.template.rename", templateID, name)

	return err
}

// RetrieveInfo retrieves information for the template.
func (vmts *VMTemplateService) RetrieveInfo(ctx context.Context, templateID int) (*resources.VMTemplate, error) {
	doc, err := vmts.retrieveInfo(ctx, "one.template.info", templateID)
	if err != nil {
		return nil, err
	}

	return resources.CreateVMTemplateFromXML(doc.Root()), nil
}

func (vmts *VMTemplateService) list(ctx context.Context, filterFlag, pageOffset,
	pageSize int) ([]*resources.VMTemplate, error) {
	resArr, err := vmts.call(ctx, "one.templatepool.info", filterFlag, pageOffset, pageSize)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromString(resArr[resultIndex].ResultString()); err != nil {
		return nil, err
	}

	elements := doc.FindElements("VMTEMPLATE_POOL/VMTEMPLATE")

	vmTemplates := make([]*resources.VMTemplate, len(elements))
	for i, e := range elements {
		vmTemplates[i] = resources.CreateVMTemplateFromXML(e)
	}

	return vmTemplates, nil
}

// ListAll retrieves information for all the VMTemplates in the pool.
func (vmts *VMTemplateService) ListAll(ctx context.Context, filter OwnershipFilter) ([]*resources.VMTemplate, error) {
	return vmts.list(ctx, int(filter), pageOffsetDefault, pageSizeDefault)
}

// ListAllForUser retrieves information for all the VMTemplates for the given user in the pool.
func (vmts *VMTemplateService) ListAllForUser(ctx context.Context,
	user resources.User) ([]*resources.VMTemplate, error) {
	userID, err := user.ID()
	if err != nil {
		return nil, err
	}

	return vmts.list(ctx, userID, pageOffsetDefault, pageSizeDefault)
}

// List retrieves information for a part of the VMTemplates in the pool with a given pagination.
func (vmts *VMTemplateService) List(ctx context.Context, pageOffset, pageSize int,
	filter OwnershipFilter) ([]*resources.VMTemplate, error) {
	return vmts.list(ctx, int(filter), (pageOffset-1)*pageSize, -pageSize)
}

// ListForUser retrieves information for a part of the VMTemplates for given user in the pool with a given pagination.
func (vmts *VMTemplateService) ListForUser(ctx context.Context, user resources.User, pageOffset,
	pageSize int) ([]*resources.VMTemplate, error) {
	userID, err := user.ID()
	if err != nil {
		return nil, err
	}

	return vmts.list(ctx, userID, (pageOffset-1)*pageSize, -pageSize)
}
