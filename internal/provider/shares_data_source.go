// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"

	"github.com/Sander0542/terraform-provider-unraid/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &SharesDataSource{}

func NewSharesDataSource() datasource.DataSource {
	return &SharesDataSource{}
}

type SharesDataSource struct {
	client *client.Client
}

type SharesDataSourceModel struct {
	Shares []ShareModel `tfsdk:"shares"`
}

type ShareModel struct {
	Name       types.String `tfsdk:"name"`
	Free       types.Int64  `tfsdk:"free"`
	Used       types.Int64  `tfsdk:"used"`
	Size       types.Int64  `tfsdk:"size"`
	Include    types.List   `tfsdk:"include"`
	Exclude    types.List   `tfsdk:"exclude"`
	Cache      types.Bool   `tfsdk:"cache"`
	NameOrig   types.String `tfsdk:"name_orig"`
	Comment    types.String `tfsdk:"comment"`
	Allocator  types.String `tfsdk:"allocator"`
	SplitLevel types.String `tfsdk:"split_level"`
	Floor      types.String `tfsdk:"floor"`
	Cow        types.String `tfsdk:"cow"`
	Color      types.String `tfsdk:"color"`
	LuksStatus types.String `tfsdk:"luks_status"`
}

func (d *SharesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shares"
}

func (d *SharesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the list of shares from the Unraid server.",
		Attributes: map[string]schema.Attribute{
			"shares": schema.ListNestedAttribute{
				MarkdownDescription: "The list of shares on the Unraid server.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The display name of the share.",
							Computed:            true,
						},
						"free": schema.Int64Attribute{
							MarkdownDescription: "Free space in kilobytes.",
							Computed:            true,
						},
						"used": schema.Int64Attribute{
							MarkdownDescription: "Used space in kilobytes.",
							Computed:            true,
						},
						"size": schema.Int64Attribute{
							MarkdownDescription: "Total size in kilobytes.",
							Computed:            true,
						},
						"include": schema.ListAttribute{
							MarkdownDescription: "Disks that are included in this share.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"exclude": schema.ListAttribute{
							MarkdownDescription: "Disks that are excluded from this share.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"cache": schema.BoolAttribute{
							MarkdownDescription: "Whether this share is cached.",
							Computed:            true,
						},
						"name_orig": schema.StringAttribute{
							MarkdownDescription: "The original name of the share.",
							Computed:            true,
						},
						"comment": schema.StringAttribute{
							MarkdownDescription: "User comment for the share.",
							Computed:            true,
						},
						"allocator": schema.StringAttribute{
							MarkdownDescription: "The allocator used for the share.",
							Computed:            true,
						},
						"split_level": schema.StringAttribute{
							MarkdownDescription: "The split level of the share.",
							Computed:            true,
						},
						"floor": schema.StringAttribute{
							MarkdownDescription: "The floor value of the share.",
							Computed:            true,
						},
						"cow": schema.StringAttribute{
							MarkdownDescription: "The COW (Copy-on-Write) setting.",
							Computed:            true,
						},
						"color": schema.StringAttribute{
							MarkdownDescription: "The color indicator of the share.",
							Computed:            true,
						},
						"luks_status": schema.StringAttribute{
							MarkdownDescription: "The LUKS encryption status.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *SharesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *SharesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, err := client.GetShares(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Unraid shares", err.Error())
		return
	}

	state := SharesDataSourceModel{
		Shares: make([]ShareModel, 0, len(data.Shares)),
	}

	for _, share := range data.Shares {
		includeList, diags := types.ListValueFrom(ctx, types.StringType, share.Include)
		resp.Diagnostics.Append(diags...)
		excludeList, diags := types.ListValueFrom(ctx, types.StringType, share.Exclude)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		state.Shares = append(state.Shares, ShareModel{
			Name:       types.StringValue(share.Name),
			Free:       types.Int64Value(share.Free),
			Used:       types.Int64Value(share.Used),
			Size:       types.Int64Value(share.Size),
			Include:    includeList,
			Exclude:    excludeList,
			Cache:      types.BoolValue(share.Cache),
			NameOrig:   types.StringValue(share.NameOrig),
			Comment:    types.StringValue(share.Comment),
			Allocator:  types.StringValue(share.Allocator),
			SplitLevel: types.StringValue(share.SplitLevel),
			Floor:      types.StringValue(share.Floor),
			Cow:        types.StringValue(share.Cow),
			Color:      types.StringValue(share.Color),
			LuksStatus: types.StringValue(share.LuksStatus),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
