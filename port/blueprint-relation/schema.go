package blueprint_relation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func BlueprintRelationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The ID of this resource, in the form `<blueprint>:<identifier>`",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The identifier of the blueprint the relation belongs to",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the relation",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"target": schema.StringAttribute{
			MarkdownDescription: "The identifier of the blueprint the relation points at",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the relation",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the relation",
			Optional:            true,
		},
		"many": schema.BoolAttribute{
			MarkdownDescription: "Whether the relation points at many entities",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"required": schema.BoolAttribute{
			MarkdownDescription: "Whether the relation is required. Port rejects a required relation whose target blueprint already has a required relation back to this one, so at most one direction of a circular pair can be required.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
	}
}

func (r *BlueprintRelationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: blueprintRelationMarkdownDescription,
		Attributes:          BlueprintRelationSchema(),
	}
}

var blueprintRelationMarkdownDescription = `

# Blueprint Relation

Docs about blueprint relations can be found [here](https://docs.port.io/context-lake/data-model/setup-blueprint/relate-blueprints).

This resource declares a single relation on a blueprint, separately from the ` + "`port_blueprint`" + ` resource
that owns the blueprint itself.

Use it when relations cannot be declared inline on ` + "`port_blueprint`" + `. The main case is two blueprints
that point at each other: Port requires a relation's target blueprint to already exist, so neither blueprint
can be created first with its relations attached. Declaring the relations as their own resources makes both
blueprints get created first, and the relations afterwards.

A relation must be managed either inline on ` + "`port_blueprint`" + ` or by this resource, never by both.
Creating a relation that already exists on the blueprint fails, rather than silently taking it over.

~> **Moving an existing relation here is destructive.** Removing a relation from a ` + "`port_blueprint`" + `
resource deletes it, which deletes the relation's values on every entity and strips any mirror property or
inherited ownership that depends on it. Re-adding it as a ` + "`port_blueprint_relation`" + ` does not bring
that data back. Only move relations that do not exist yet.

## Example Usage

Two blueprints that relate to each other. Neither blueprint references the other, so there is no cycle,
and both are created before either relation is written:

` + "```hcl" + `

resource "port_blueprint" "entra_id_user" {
  title      = "Entra ID User"
  identifier = "entra-id-user"
  icon       = "User"
}

resource "port_blueprint" "entra_id_group" {
  title      = "Entra ID Group"
  identifier = "entra-id-group"
  icon       = "Group"
}

resource "port_blueprint_relation" "entra_id_user_to_group" {
  blueprint  = port_blueprint.entra_id_user.identifier
  identifier = "groups"
  target     = port_blueprint.entra_id_group.identifier
  title      = "Groups"
  many       = true
}

resource "port_blueprint_relation" "entra_id_group_to_user" {
  blueprint  = port_blueprint.entra_id_group.identifier
  identifier = "members"
  target     = port_blueprint.entra_id_user.identifier
  title      = "Members"
  many       = true
}

` + "```" + `

## Import

Blueprint relations are imported using ` + "`<blueprint>:<identifier>`" + `:

` + "```shell" + `
terraform import port_blueprint_relation.members "entra-id-group:members"
` + "```" + `
`
