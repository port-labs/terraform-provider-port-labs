package scorecard_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/port/scorecard"
)

func memberSpecSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"filter": schema.SingleNestedAttribute{
			MarkdownDescription: "An optional set of conditions to filter entities evaluated by this member scorecard.",
			Optional:            true,
			Attributes:          querySchema(),
		},
		"rules": schema.ListNestedAttribute{
			MarkdownDescription: "The rules that define this member scorecard.",
			Required:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: scorecard.RuleSchema(),
			},
		},
	}
}

func querySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"combinator": schema.StringAttribute{
			MarkdownDescription: "The combinator of the query.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("and", "or"),
			},
		},
		"conditions": schema.ListAttribute{
			MarkdownDescription: "The conditions of the query. Each condition object should be encoded to a string.",
			Required:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
		},
	}
}

func (r *ScorecardGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: resourceMarkdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "A unique identifier for the scorecard group.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The title of the scorecard group (applied to all member scorecards).",
				Required:            true,
			},
			"levels": schema.ListNestedAttribute{
				MarkdownDescription: "The available levels of the scorecard group, shared by all members.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: scorecard.LevelSchema(),
				},
			},
			"blueprints": schema.SetAttribute{
				MarkdownDescription: "Blueprint identifiers that share the same rules (and optional filters). Use this for shared-rules mode. Conflicts with `scorecards`.",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ConflictsWith(path.MatchRoot("scorecards")),
					setvalidator.AlsoRequires(path.MatchRoot("rules")),
				},
			},
			"rules": schema.ListNestedAttribute{
				MarkdownDescription: "The rules applied to every blueprint in shared-rules mode. Conflicts with `scorecards`.",
				Optional:            true,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("scorecards")),
					listvalidator.AlsoRequires(path.MatchRoot("blueprints")),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: scorecard.RuleSchema(),
				},
			},
			"filters": schema.MapNestedAttribute{
				MarkdownDescription: "Optional filters per blueprint in shared-rules mode, keyed by blueprint identifier. Conflicts with `scorecards`.",
				Optional:            true,
				Validators: []validator.Map{
					mapvalidator.ConflictsWith(path.MatchRoot("scorecards")),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: querySchema(),
				},
			},
			"scorecards": schema.MapNestedAttribute{
				MarkdownDescription: "Map of blueprint identifier to member scorecard filter/rules. Use this for per-blueprint mode. Conflicts with `blueprints`, `rules`, and `filters`.",
				Optional:            true,
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
					mapvalidator.ConflictsWith(
						path.MatchRoot("blueprints"),
						path.MatchRoot("rules"),
						path.MatchRoot("filters"),
					),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: memberSpecSchema(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The creation date of the scorecard group.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "The creator of the scorecard group.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The last update date of the scorecard group.",
				Computed:            true,
			},
			"updated_by": schema.StringAttribute{
				MarkdownDescription: "The last updater of the scorecard group.",
				Computed:            true,
			},
		},
	}
}

var resourceMarkdownDescription = `
# Scorecard Group

This resource allows you to manage a scorecard group that creates scorecards across multiple blueprints.

~> **Note:** Scorecard groups are currently protected by the ` + "`SCORECARD_GROUPS`" + ` organization feature flag and are not available in all Port organizations. If you need access, contact your Port account team.

A scorecard group can be configured in one of two modes:


- **Shared rules mode** — set ` + "`blueprints`" + `, ` + "`rules`" + `, and optionally ` + "`filters`" + ` to apply the same rules to multiple blueprints.
- **Per-blueprint mode** — set ` + "`scorecards`" + ` to define different filter/rules per blueprint.

See the [Port documentation](https://docs.getport.io/governance/standards-and-compliance/manage-scorecards/) for more information about scorecards.

## Example Usage (shared rules)

` + "```hcl" + `

resource "port_scorecard_group" "readiness" {
  identifier = "production-readiness"
  title      = "Production Readiness"
  blueprints = [port_blueprint.microservice.identifier]
  rules = [{
    identifier = "has-owner"
    title      = "Has Owner"
    level      = "Gold"
    query = {
      combinator = "and"
      conditions = [jsonencode({
        property = "$team"
        operator = "isNotEmpty"
      })]
    }
  }]
}

` + "```" + `

## Example Usage (per blueprint)

` + "```hcl" + `

resource "port_scorecard_group" "readiness" {
  identifier = "production-readiness"
  title      = "Production Readiness"
  scorecards = {
    (port_blueprint.microservice.identifier) = {
      rules = [{
        identifier = "has-owner"
        title      = "Has Owner"
        level      = "Gold"
        query = {
          combinator = "and"
          conditions = [jsonencode({
            property = "$team"
            operator = "isNotEmpty"
          })]
        }
      }]
    }
  }
}

` + "```"
