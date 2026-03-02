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
	t.Helper()

	ctx := t.Context()
	schemaResp := testSchema()

	partnerElems := make([]types.String, len(partners))
	for i, p := range partners {
		partnerElems[i] = types.StringValue(p)
	}

	childElems := make([]types.String, len(children))
	for i, c := range children {
		childElems[i] = types.StringValue(c)
	}

	partnersSet, diags := types.SetValueFrom(ctx, types.StringType, partnerElems)
	if diags.HasError() {
		t.Fatalf("failed to create partners set: %v", diags)
	}

	childrenSet, diags := types.SetValueFrom(ctx, types.StringType, childElems)
	if diags.HasError() {
		t.Fatalf("failed to create children set: %v", diags)
	}

	model := ResourceModel{
		ID:       types.StringValue("union-1"),
		Partners: partnersSet,
		Children: childrenSet,
		Marriage: types.ObjectNull(event.EventModelAttributeTypes()),
		Divorce:  types.ObjectNull(event.EventModelAttributeTypes()),
	}

	objVal, diags := types.ObjectValueFrom(ctx, schemaResp.Schema.Type().(types.ObjectType).AttrTypes, model)
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
}
