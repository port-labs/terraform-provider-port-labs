package cli

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

// ResourceExistsError is returned when creating a resource whose object already exists in Port, and the
// resource is not allowed to take ownership of it
type ResourceExistsError struct {
	ResourceKind consts.ResourceKind
	Identifier   string
	// Blueprint is the blueprint the object belongs to, for objects that live on a blueprint such as relations
	Blueprint string
}

func (e *ResourceExistsError) Summary() string {
	return fmt.Sprintf("%s already exists", e.ResourceKind.Noun())
}

func (e *ResourceExistsError) Error() string {
	location := "in Port"
	importID := e.Identifier
	remedy := ""
	if e.Blueprint != "" {
		location = fmt.Sprintf("on blueprint %q", e.Blueprint)
		importID = fmt.Sprintf("%s:%s", e.Blueprint, e.Identifier)
		remedy = "Remove it from the blueprint's relations block, "
	}
	return fmt.Sprintf("%s %q already exists %s. %sImport it with `terraform import <address> %q`, "+
		"or set `adopt_existing = true` (or add %q to the provider's `adopt_existing_resources`) to take ownership of it on creation.",
		e.ResourceKind.Noun(), e.Identifier, location, remedy, importID, e.ResourceKind)
}

// IdentifierTakenError is returned when Port rejects creating an object because its identifier is already
// taken, e.g. by an object Port created asynchronously after the resource checked that it did not exist
type IdentifierTakenError struct {
	Identifier string
	Response   string
}

func (e *IdentifierTakenError) Error() string {
	return fmt.Sprintf("identifier %q is already taken, got: %s", e.Identifier, e.Response)
}

// AdoptionConflictError is returned when an object that already exists in Port cannot be taken over, because a
// field that cannot be changed after creation differs from the configuration
type AdoptionConflictError struct {
	ResourceKind consts.ResourceKind
	Identifier   string
	Field        string
	Existing     string
	Configured   string
}

func (e *AdoptionConflictError) Summary() string {
	return fmt.Sprintf("cannot take over existing %s", e.ResourceKind.Noun())
}

func (e *AdoptionConflictError) Error() string {
	return fmt.Sprintf("%s %q already exists in Port with `%s` = %q, but the configuration sets %q. `%s` cannot be changed "+
		"after creation, so the existing %s cannot be taken over. Set `%s` to %q, or remove the existing %s from Port.",
		e.ResourceKind.Noun(), e.Identifier, e.Field, e.Existing, e.Configured, e.Field,
		e.ResourceKind.Noun(), e.Field, e.Existing, e.ResourceKind.Noun())
}

// AdoptionPolicy holds the resource kinds that may take ownership of an object that already exists in Port,
// as configured by the provider's adopt_existing_resources
type AdoptionPolicy map[consts.ResourceKind]bool

func NewAdoptionPolicy(kinds ...consts.ResourceKind) AdoptionPolicy {
	policy := make(AdoptionPolicy, len(kinds))
	for _, kind := range kinds {
		policy[kind] = true
	}
	return policy
}

func (p AdoptionPolicy) Allows(kind consts.ResourceKind) bool {
	return p[kind]
}

// AdoptExisting reports whether a resource of the given kind may take ownership of an object that already exists
// in Port. The resource's own adopt_existing setting wins; when unset, it is allowed only if the kind is listed
// in the provider's adopt_existing_resources.
func (c *PortClient) AdoptExisting(kind consts.ResourceKind, resourceSetting types.Bool) bool {
	if resourceSetting.IsNull() || resourceSetting.IsUnknown() {
		return c.AdoptExistingResources.Allows(kind)
	}
	return resourceSetting.ValueBool()
}
