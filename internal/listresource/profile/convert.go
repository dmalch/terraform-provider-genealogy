// Package profilelist implements the `geni_profile` list resource for
// Terraform 1.14's `terraform query`. The package translates pages of
// managed-profile API responses into list.ListResult values whose Identity
// matches the managed resource's identity schema, so practitioners can pipe
// the output directly into an `import { identity = ... }` block.
package profilelist

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	profileresource "github.com/dmalch/terraform-provider-genealogy/internal/resource/profile"
)

// displayNameFor produces a human-readable label for a profile in query output.
// The Geni API returns both a localized `names` map and flat top-level
// name fields; we prefer the localized en-US entry when present and fall
// back to the flat fields otherwise. When nothing is populated the bare ID
// is returned so the result is still visually identifiable.
func displayNameFor(p *geni.ProfileResponse) string {
	first, last := canonicalFirstLast(p)
	switch {
	case first != "" && last != "":
		return fmt.Sprintf("%s %s (%s)", first, last, p.Id)
	case first != "":
		return fmt.Sprintf("%s (%s)", first, p.Id)
	case last != "":
		return fmt.Sprintf("%s (%s)", last, p.Id)
	default:
		return p.Id
	}
}

func canonicalFirstLast(p *geni.ProfileResponse) (string, string) {
	if name, ok := p.Names["en-US"]; ok {
		return stringFromPtr(name.FirstName), stringFromPtr(name.LastName)
	}
	return stringFromPtr(p.FirstName), stringFromPtr(p.LastName)
}

func stringFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// buildProfileListResult turns one API response into a list.ListResult whose
// Identity carries the profile ID under the managed resource's identity
// schema. When req.IncludeResource is true the Resource field is populated
// via the managed resource's own ValueFrom translator — the same one used
// by the managed Read flow — so list output round-trips through `import`.
//
// Returning false signals the caller (paginate) to stop iteration; in that
// case the result's Diagnostics already carry the offending error.
func buildProfileListResult(
	ctx context.Context,
	resp *geni.ProfileResponse,
	req list.ListRequest,
) (list.ListResult, bool) {
	result := req.NewListResult(ctx)

	identity := profileresource.ResourceIdentityModel{
		ID: types.StringValue(resp.Id),
	}
	diags := result.Identity.Set(ctx, identity)
	result.Diagnostics.Append(diags...)
	if result.Diagnostics.HasError() {
		return result, false
	}

	result.DisplayName = displayNameFor(resp)

	if req.IncludeResource {
		// Seed every collection field with a typed null so ValueFrom can
		// overwrite the fields the API does return without leaving any
		// untouched attribute carrying a type-less zero `types.Set{}` /
		// `types.Map{}` / `types.Object{}`, which the framework rejects with
		// a MISSING TYPE conversion error.
		model := profileresource.NewEmptyResourceModel()

		diags = profileresource.ValueFrom(ctx, resp, &model)
		result.Diagnostics.Append(diags...)
		if result.Diagnostics.HasError() {
			return result, false
		}

		diags = result.Resource.Set(ctx, model)
		result.Diagnostics.Append(diags...)
		if result.Diagnostics.HasError() {
			return result, false
		}
	}

	return result, true
}
