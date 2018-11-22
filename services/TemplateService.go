package services

import (
	"context"

	"github.com/beevik/etree"

	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/requests"
	"github.com/onego-project/onego/resources"
)

// TemplateService structure to manage OpenNebula Template.
type TemplateService struct {
	Service
}

// Allocate allocates a new template in OpenNebula.
func (ts *TemplateService) Allocate(ctx context.Context, blueprint blueprint.Interface) (*resources.Template, error) {
	blueprintText, err := blueprint.Render()
	if err != nil {
		return nil, err
	}

	resArr, err := ts.call(ctx, "one.template.allocate", blueprintText)
	if err != nil {
		return nil, err
	}

	return ts.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Clone clones an existing virtual machine template.
func (ts *TemplateService) Clone(ctx context.Context, template resources.Template, name string,
	withImages bool) (*resources.Template, error) {
	templateID, err := template.ID()
	if err != nil {
		return nil, err
	}

	resArr, err := ts.call(ctx, "one.template.clone", templateID, name, withImages)
	if err != nil {
		return nil, err
	}

	return ts.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Delete deletes the given template from the pool.
func (ts *TemplateService) Delete(ctx context.Context, template resources.Template, withImages bool) error {
	templateID, err := template.ID()
	if err != nil {
		return err
	}

	_, err = ts.call(ctx, "one.template.delete", templateID, withImages)

	return err
}

// Update merges or replaces the template contents.
func (ts *TemplateService) Update(ctx context.Context, template resources.Template, blueprint blueprint.Interface,
	updateType UpdateType) (*resources.Template, error) {
	templateID, err := template.ID()
	if err != nil {
		return nil, err
	}

	blueprintText, err := blueprint.Render()
	if err != nil {
		return nil, err
	}

	resArr, err := ts.call(ctx, "one.template.update", templateID, blueprintText, updateType)
	if err != nil {
		return nil, err
	}

	return ts.RetrieveInfo(ctx, int(resArr[resultIndex].ResultInt()))
}

// Chmod changes the permission bits of a template.
func (ts *TemplateService) Chmod(ctx context.Context, template resources.Template,
	request requests.PermissionRequest) error {
	templateID, err := template.ID()
	if err != nil {
		return err
	}

	return ts.chmod(ctx, "one.template.chmod", templateID, request)
}

// Chown changes the ownership of a template.
func (ts *TemplateService) Chown(ctx context.Context, template resources.Template,
	request requests.OwnershipRequest) error {
	templateID, err := template.ID()
	if err != nil {
		return err
	}

	return ts.chown(ctx, "one.template.chown", templateID, request)
}

// Rename renames a template.
func (ts *TemplateService) Rename(ctx context.Context, template resources.Template, name string) error {
	templateID, err := template.ID()
	if err != nil {
		return err
	}

	_, err = ts.call(ctx, "one.template.rename", templateID, name)

	return err
}

// RetrieveInfo retrieves information for the template.
func (ts *TemplateService) RetrieveInfo(ctx context.Context, templateID int) (*resources.Template, error) {
	doc, err := ts.retrieveInfo(ctx, "one.template.info", templateID)
	if err != nil {
		return nil, err
	}

	return resources.CreateTemplateFromXML(doc.Root()), nil
}

func (ts *TemplateService) list(ctx context.Context, filterFlag, pageOffset,
	pageSize int) ([]*resources.Template, error) {
	resArr, err := ts.call(ctx, "one.templatepool.info", filterFlag, pageOffset, pageSize)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err = doc.ReadFromString(resArr[resultIndex].ResultString()); err != nil {
		return nil, err
	}

	elements := doc.FindElements("VMTEMPLATE_POOL/VMTEMPLATE")

	vmTemplates := make([]*resources.Template, len(elements))
	for i, e := range elements {
		vmTemplates[i] = resources.CreateTemplateFromXML(e)
	}

	return vmTemplates, nil
}

// ListAll retrieves information for all the VMTemplates in the pool.
func (ts *TemplateService) ListAll(ctx context.Context, filter OwnershipFilter) ([]*resources.Template, error) {
	return ts.list(ctx, int(filter), pageOffsetDefault, pageSizeDefault)
}

// ListAllForUser retrieves information for all the VMTemplates for the given user in the pool.
func (ts *TemplateService) ListAllForUser(ctx context.Context,
	user resources.User) ([]*resources.Template, error) {
	userID, err := user.ID()
	if err != nil {
		return nil, err
	}

	return ts.list(ctx, userID, pageOffsetDefault, pageSizeDefault)
}

// List retrieves information for a part of the VMTemplates in the pool with a given pagination.
func (ts *TemplateService) List(ctx context.Context, pageOffset, pageSize int,
	filter OwnershipFilter) ([]*resources.Template, error) {
	return ts.list(ctx, int(filter), (pageOffset-1)*pageSize, -pageSize)
}

// ListForUser retrieves information for a part of the VMTemplates for given user in the pool with a given pagination.
func (ts *TemplateService) ListForUser(ctx context.Context, user resources.User, pageOffset,
	pageSize int) ([]*resources.Template, error) {
	userID, err := user.ID()
	if err != nil {
		return nil, err
	}

	return ts.list(ctx, userID, (pageOffset-1)*pageSize, -pageSize)
}
