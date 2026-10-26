package consts

// ResourceKind is the type name of a resource this provider manages, as written in Terraform configuration
type ResourceKind string

const (
	ResourceKindPage              ResourceKind = "port_page"
	ResourceKindBlueprint         ResourceKind = "port_blueprint"
	ResourceKindBlueprintRelation ResourceKind = "port_blueprint_relation"
	ResourceKindIntegration       ResourceKind = "port_integration"
)

// AdoptableResourceKinds are the resource kinds that can take ownership of an object that already exists in Port,
// as accepted by the provider's adopt_existing_resources
var AdoptableResourceKinds = []ResourceKind{
	ResourceKindPage,
	ResourceKindBlueprint,
	ResourceKindBlueprintRelation,
	ResourceKindIntegration,
}

// Noun is how the resource kind's objects are called in messages
func (k ResourceKind) Noun() string {
	switch k {
	case ResourceKindPage:
		return "page"
	case ResourceKindBlueprint:
		return "blueprint"
	case ResourceKindBlueprintRelation:
		return "relation"
	case ResourceKindIntegration:
		return "integration"
	}
	return string(k)
}
