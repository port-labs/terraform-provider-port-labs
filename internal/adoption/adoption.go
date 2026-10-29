// Package adoption decides, in one place for every resource, what creating a resource does about an object that
// already exists in Port: fail, or take ownership of it and say so.
package adoption

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

// Phase is when the existing object was found: while planning, before anything is written, or while applying,
// including objects Port created after the plan (e.g. a catalog page created alongside its blueprint)
type Phase int

const (
	AtPlan Phase = iota
	AtApply
)

// TakeOver reports whether the resource may take ownership of the existing object. It adds to diags the error to
// fail with when it may not (taking over is not allowed, or conflict is set), and otherwise a warning naming the
// object being taken over, so taking over is never silent.
func TakeOver(phase Phase, client *cli.PortClient, existing *cli.ResourceExistsError, adoptExisting types.Bool, conflict *cli.AdoptionConflictError, diags *diag.Diagnostics) bool {
	if !client.AdoptExisting(existing.ResourceKind, adoptExisting) {
		diags.AddAttributeError(identifierPath(existing.ResourceKind), existing.Summary(), existing.Error())
		return false
	}
	if conflict != nil {
		diags.AddAttributeError(path.Root(conflict.Field), conflict.Summary(), conflict.Error())
		return false
	}

	noun := existing.ResourceKind.Noun()
	if phase == AtPlan {
		diags.AddAttributeWarning(path.Root("adopt_existing"), fmt.Sprintf("%s will be taken over", noun),
			fmt.Sprintf("%s %q already exists %s. Applying will take ownership of it and overwrite it to match the configuration.",
				noun, existing.Identifier, location(existing)))
		return true
	}
	diags.AddAttributeWarning(path.Root("adopt_existing"), fmt.Sprintf("%s was taken over", noun),
		fmt.Sprintf("%s %q already existed %s, created outside Terraform. Terraform took ownership of it and overwrote it to match the configuration.",
			noun, existing.Identifier, location(existing)))
	return true
}

func identifierPath(kind consts.ResourceKind) path.Path {
	if kind == consts.ResourceKindIntegration {
		return path.Root("installation_id")
	}
	return path.Root("identifier")
}

func location(existing *cli.ResourceExistsError) string {
	if existing.Blueprint != "" {
		return fmt.Sprintf("on blueprint %q", existing.Blueprint)
	}
	return "in Port"
}
