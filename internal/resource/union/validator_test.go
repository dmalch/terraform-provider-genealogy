package union

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

func testSchema() resource.SchemaResponse {
	r := &Resource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	return *resp
}

func validatorTestConfig(t *testing.T, partners []string, children []string) tfsdk.Config {
	return validatorTestConfigFull(t, partners, children, nil, nil)
}

func validatorTestConfigFull(t *testing.T, partners, children, fosterChildren, adoptedChildren []string) tfsdk.Config {
	t.Helper()

	ctx := t.Context()
	schemaResp := testSchema()

	toSet := func(ids []string, field string) types.Set {
		elems := make([]types.String, len(ids))
		for i, id := range ids {
			elems[i] = types.StringValue(id)
		}
		set, diags := types.SetValueFrom(ctx, types.StringType, elems)
		if diags.HasError() {
			t.Fatalf("failed to create %s set: %v", field, diags)
		}
		return set
	}

	model := ResourceModel{
		ID:              types.StringValue("union-1"),
		Partners:        toSet(partners, "partners"),
		Children:        toSet(children, "children"),
		FosterChildren:  toSet(fosterChildren, "foster_children"),
		AdoptedChildren: toSet(adoptedChildren, "adopted_children"),
		Marriage:        types.ObjectNull(event.EventModelAttributeTypes()),
		Divorce:         types.ObjectNull(event.EventModelAttributeTypes()),
	}

	objType, ok := schemaResp.Schema.Type().(types.ObjectType)
	if !ok {
		t.Fatalf("schema type is not an object: %T", schemaResp.Schema.Type())
	}
	objVal, diags := types.ObjectValueFrom(ctx, objType.AttrTypes, model)
	if diags.HasError() {
		t.Fatalf("failed to create object value: %v", diags)
	}

	tfVal, err := objVal.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("failed to convert to terraform value: %v", err)
	}

	return tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw:    tfVal,
	}
}

func TestValidateConfig(t *testing.T) {
	t.Run("valid config with 2 partners", func(t *testing.T) {
		RegisterTestingT(t)

		r := &Resource{}
		req := resource.ValidateConfigRequest{
			Config: validatorTestConfig(t, []string{"profile-1", "profile-2"}, []string{}),
		}
		resp := &resource.ValidateConfigResponse{}

		r.ValidateConfig(t.Context(), req, resp)

		Expect(resp.Diagnostics.HasError()).To(BeFalse())
	})

	t.Run("too many partners", func(t *testing.T) {
		RegisterTestingT(t)

		r := &Resource{}
		req := resource.ValidateConfigRequest{
			Config: validatorTestConfig(t, []string{"profile-1", "profile-2", "profile-3"}, []string{}),
		}
		resp := &resource.ValidateConfigResponse{}

		r.ValidateConfig(t.Context(), req, resp)

		Expect(resp.Diagnostics.HasError()).To(BeTrue())
		Expect(resp.Diagnostics.Errors()).To(HaveLen(1))
		Expect(resp.Diagnostics.Errors()[0].Summary()).To(Equal("Too Many Partners"))
	})

	t.Run("insufficient profiles with 1 partner and 0 children", func(t *testing.T) {
		RegisterTestingT(t)

		r := &Resource{}
		req := resource.ValidateConfigRequest{
			Config: validatorTestConfig(t, []string{"profile-1"}, []string{}),
		}
		resp := &resource.ValidateConfigResponse{}

		r.ValidateConfig(t.Context(), req, resp)

		Expect(resp.Diagnostics.HasError()).To(BeTrue())
		Expect(resp.Diagnostics.Errors()[0].Summary()).To(Equal("Insufficient Attribute Configuration"))
	})

	t.Run("children compensate for missing partners", func(t *testing.T) {
		RegisterTestingT(t)

		r := &Resource{}
		req := resource.ValidateConfigRequest{
			Config: validatorTestConfig(t, []string{}, []string{"profile-1", "profile-2"}),
		}
		resp := &resource.ValidateConfigResponse{}

		r.ValidateConfig(t.Context(), req, resp)

		Expect(resp.Diagnostics.HasError()).To(BeFalse())
	})

	t.Run("mixed valid with 1 partner and 1 child", func(t *testing.T) {
		RegisterTestingT(t)

		r := &Resource{}
		req := resource.ValidateConfigRequest{
			Config: validatorTestConfig(t, []string{"profile-1"}, []string{"profile-2"}),
		}
		resp := &resource.ValidateConfigResponse{}

		r.ValidateConfig(t.Context(), req, resp)

		Expect(resp.Diagnostics.HasError()).To(BeFalse())
	})

	t.Run("children and foster_children must be disjoint", func(t *testing.T) {
		RegisterTestingT(t)

		r := &Resource{}
		req := resource.ValidateConfigRequest{
			Config: validatorTestConfigFull(t,
				[]string{"profile-1", "profile-2"},
				[]string{"profile-3"},
				[]string{"profile-3"}, // overlap with children
				nil,
			),
		}
		resp := &resource.ValidateConfigResponse{}

		r.ValidateConfig(t.Context(), req, resp)

		Expect(resp.Diagnostics.HasError()).To(BeTrue())
		Expect(resp.Diagnostics.Errors()[0].Summary()).To(Equal("Overlapping Child Sets"))
	})

	t.Run("children and adopted_children must be disjoint", func(t *testing.T) {
		RegisterTestingT(t)

		r := &Resource{}
		req := resource.ValidateConfigRequest{
			Config: validatorTestConfigFull(t,
				[]string{"profile-1", "profile-2"},
				[]string{"profile-3"},
				nil,
				[]string{"profile-3"}, // overlap with children
			),
		}
		resp := &resource.ValidateConfigResponse{}

		r.ValidateConfig(t.Context(), req, resp)

		Expect(resp.Diagnostics.HasError()).To(BeTrue())
		Expect(resp.Diagnostics.Errors()[0].Summary()).To(Equal("Overlapping Child Sets"))
	})

	t.Run("foster_children and adopted_children must be disjoint", func(t *testing.T) {
		RegisterTestingT(t)

		r := &Resource{}
		req := resource.ValidateConfigRequest{
			Config: validatorTestConfigFull(t,
				[]string{"profile-1", "profile-2"},
				nil,
				[]string{"profile-3"},
				[]string{"profile-3"}, // overlap with foster
			),
		}
		resp := &resource.ValidateConfigResponse{}

		r.ValidateConfig(t.Context(), req, resp)

		Expect(resp.Diagnostics.HasError()).To(BeTrue())
		Expect(resp.Diagnostics.Errors()[0].Summary()).To(Equal("Overlapping Child Sets"))
	})

	t.Run("all three child sets disjoint is valid", func(t *testing.T) {
		RegisterTestingT(t)

		r := &Resource{}
		req := resource.ValidateConfigRequest{
			Config: validatorTestConfigFull(t,
				[]string{"profile-1", "profile-2"},
				[]string{"profile-3"},
				[]string{"profile-4"},
				[]string{"profile-5"},
			),
		}
		resp := &resource.ValidateConfigResponse{}

		r.ValidateConfig(t.Context(), req, resp)

		Expect(resp.Diagnostics.HasError()).To(BeFalse())
	})

	t.Run("min-2-profiles counts across all child sets", func(t *testing.T) {
		RegisterTestingT(t)

		r := &Resource{}
		// 0 partners + 0 children + 1 foster + 1 adopted = 2 profiles → valid
		req := resource.ValidateConfigRequest{
			Config: validatorTestConfigFull(t,
				nil,
				nil,
				[]string{"profile-1"},
				[]string{"profile-2"},
			),
		}
		resp := &resource.ValidateConfigResponse{}

		r.ValidateConfig(t.Context(), req, resp)

		Expect(resp.Diagnostics.HasError()).To(BeFalse())
	})
}
