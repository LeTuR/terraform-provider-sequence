// Copyright (c) Arthur Cesaré-Herriau
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource               = &NumberResource{}
	_ resource.ResourceWithModifyPlan = &NumberResource{}
)

func NewNumberResource() resource.Resource {
	return &NumberResource{}
}

type NumberResource struct{}

type NumberResourceModel struct {
	Start     types.Int64  `tfsdk:"start"`
	Width     types.Int64  `tfsdk:"width"`
	Prefix    types.String `tfsdk:"prefix"`
	Suffix    types.String `tfsdk:"suffix"`
	Keepers   types.Map    `tfsdk:"keepers"`
	Number    types.Int64  `tfsdk:"number"`
	Formatted types.String `tfsdk:"formatted"`
	ID        types.String `tfsdk:"id"`
}

func (r *NumberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_number"
}

func (r *NumberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a zero-padded sequential number. The number increments by 1 each time the `keepers` map changes.",
		Attributes: map[string]schema.Attribute{
			"start": schema.Int64Attribute{
				MarkdownDescription: "First value emitted on resource creation. Defaults to 1.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"width": schema.Int64Attribute{
				MarkdownDescription: "Zero-padding width for the formatted output. `width = 3` → `\"001\"`. Defaults to 3.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(3),
				Validators: []validator.Int64{
					nonNegativeInt64Validator{},
				},
			},
			"prefix": schema.StringAttribute{
				MarkdownDescription: "String prepended to the padded number. Defaults to empty.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"suffix": schema.StringAttribute{
				MarkdownDescription: "String appended to the padded number. Defaults to empty.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"keepers": schema.MapAttribute{
				MarkdownDescription: "Arbitrary map of string values. Changing any value increments `number` by 1.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"number": schema.Int64Attribute{
				MarkdownDescription: "The current sequence number.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"formatted": schema.StringAttribute{
				MarkdownDescription: "`prefix + zero-pad(number, width) + suffix`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Same as `formatted`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *NumberResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var state, plan NumberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	keepersChanged := !plan.Keepers.Equal(state.Keepers)

	var number int64
	if keepersChanged {
		number = state.Number.ValueInt64() + 1
	} else {
		number = state.Number.ValueInt64()
	}

	plan.Number = types.Int64Value(number)
	formatted := formatNumber(plan.Prefix.ValueString(), number, plan.Width.ValueInt64(), plan.Suffix.ValueString())
	plan.Formatted = types.StringValue(formatted)
	plan.ID = types.StringValue(formatted)

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func (r *NumberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NumberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	number := data.Start.ValueInt64()
	data.Number = types.Int64Value(number)
	formatted := formatNumber(data.Prefix.ValueString(), number, data.Width.ValueInt64(), data.Suffix.ValueString())
	data.Formatted = types.StringValue(formatted)
	data.ID = types.StringValue(formatted)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NumberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NumberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NumberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NumberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NumberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func formatNumber(prefix string, number int64, width int64, suffix string) string {
	if width <= 0 {
		return fmt.Sprintf("%s%d%s", prefix, number, suffix)
	}
	return fmt.Sprintf("%s%0*d%s", prefix, width, number, suffix)
}

type nonNegativeInt64Validator struct{}

func (v nonNegativeInt64Validator) Description(_ context.Context) string {
	return "value must be >= 0"
}

func (v nonNegativeInt64Validator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v nonNegativeInt64Validator) ValidateInt64(_ context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if req.ConfigValue.ValueInt64() < 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid value",
			fmt.Sprintf("Attribute %s must be >= 0, got: %d", req.Path, req.ConfigValue.ValueInt64()),
		)
	}
}
