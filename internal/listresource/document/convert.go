// Package documentlist implements the `geni_document` list resource for
// Terraform 1.14's `terraform query`. The package translates pages of
// uploaded-document API responses into list.ListResult values whose Identity
// matches the managed resource's identity schema.
package documentlist

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	documentresource "github.com/dmalch/terraform-provider-genealogy/internal/resource/document"
)

// displayNameFor produces a human-readable label for a document in query
// output. The Title is the obvious choice; we fall back to the bare ID when
// it is absent so the result remains visually identifiable.
func displayNameFor(d *geni.DocumentResponse) string {
	if d.Title == "" {
		return d.Id
	}
	return fmt.Sprintf("%s (%s)", d.Title, d.Id)
}

// buildDocumentListResult turns one API response into a list.ListResult whose
// Identity carries the document ID under the managed resource's identity
// schema. When req.IncludeResource is true the Resource field is populated
// via the managed resource's own ValueFrom translator — the same one used
// by the managed Read flow — so list output round-trips through `import`.
func buildDocumentListResult(
	ctx context.Context,
	resp *geni.DocumentResponse,
	req list.ListRequest,
) (list.ListResult, bool) {
	result := req.NewListResult(ctx)

	identity := documentresource.ResourceIdentityModel{
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
		// untouched attribute carrying a type-less zero value (see the
		// equivalent comment in profilelist for the failure mode).
		model := documentresource.NewEmptyResourceModel()

		diags = documentresource.ValueFrom(ctx, resp, &model)
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
