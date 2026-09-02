package oauth_app

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const maxRedirectURIsPerApp = 5

func OAuthAppSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The ID of the OAuth app registration",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The display name of the OAuth app registration",
			Required:            true,
		},
		"redirect_uris": schema.ListAttribute{
			MarkdownDescription: fmt.Sprintf(
				"Between 1 and %d exact, non-wildcard redirect URIs for the OAuth app. Must be absolute URIs.",
				maxRedirectURIsPerApp,
			),
			Required:    true,
			ElementType: types.StringType,
		},
		"client_id": schema.StringAttribute{
			MarkdownDescription: "The OAuth client ID assigned to this app registration",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"client_secret": schema.StringAttribute{
			MarkdownDescription: "The OAuth client secret. Returned only when the app is created and not available on subsequent reads.",
			Optional:            true,
			Computed:            true,
			Sensitive:           true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the OAuth app registration",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the OAuth app registration",
			Computed:            true,
		},
	}
}

func (r *OAuthAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: ResourceMarkdownDescription,
		Attributes:          OAuthAppSchema(),
	}
}

var ResourceMarkdownDescription = `
# OAuth App resource

OAuth app registrations allow you to register confidential OAuth clients for use with Port integrations and MCP connectors.

The ` + "`client_secret`" + ` is returned only when the app is created. Port does not expose it on subsequent API reads, so Terraform retains the value from the initial create in state.

## Example Usage

` + "```hcl" + `
resource "port_oauth_app" "mcp_connector" {
  name = "My MCP Connector"
  redirect_uris = [
    "https://api.port.io/v1/mcp/oauth2/callback",
  ]
}
` + "\n```"

var safeRedirectURIChars = regexp.MustCompile(`^[a-zA-Z0-9\-._~:/?[\]@!$&'()+,;=%]+$`)

func validateRedirectURI(redirectURI string) error {
	if !safeRedirectURIChars.MatchString(redirectURI) {
		return fmt.Errorf("redirect_uris contains a URI with characters not permitted in a URI")
	}

	if strings.Contains(redirectURI, "*") {
		return fmt.Errorf("redirect_uris must contain absolute, exact (non-wildcard) URIs")
	}

	parsed, err := url.Parse(redirectURI)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("redirect_uris must contain absolute, exact (non-wildcard) URIs")
	}

	isLocalhost := parsed.Hostname() == "localhost" || parsed.Hostname() == "127.0.0.1"
	if parsed.Scheme != "https" && !(parsed.Scheme == "http" && isLocalhost) {
		return fmt.Errorf("redirect_uris must use https, or http only for localhost")
	}

	return nil
}

func validateRedirectURIs(redirectURIs []string) error {
	if len(redirectURIs) < 1 {
		return fmt.Errorf("redirect_uris must contain at least 1 URI")
	}

	if len(redirectURIs) > maxRedirectURIsPerApp {
		return fmt.Errorf("redirect_uris must contain at most %d URIs", maxRedirectURIsPerApp)
	}

	seen := make(map[string]struct{}, len(redirectURIs))
	for _, redirectURI := range redirectURIs {
		if err := validateRedirectURI(redirectURI); err != nil {
			return err
		}

		if _, exists := seen[redirectURI]; exists {
			return fmt.Errorf("redirect_uris must contain unique URIs")
		}
		seen[redirectURI] = struct{}{}
	}

	return nil
}

func redirectURIsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sortedA := slices.Clone(a)
	sortedB := slices.Clone(b)
	slices.Sort(sortedA)
	slices.Sort(sortedB)

	return slices.Equal(sortedA, sortedB)
}
