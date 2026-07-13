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

var _ datasource.DataSource = &VersionDataSource{}

func NewVersionDataSource() datasource.DataSource {
	return &VersionDataSource{}
}

type VersionDataSource struct {
	client *client.Client
}

type VersionDataSourceModel struct {
	Unraid types.String `tfsdk:"unraid"`
	Api    types.String `tfsdk:"api"`
	Kernel types.String `tfsdk:"kernel"`
}

func (d *VersionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_version"
}

func (d *VersionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves version information from the Unraid server.",
		Attributes: map[string]schema.Attribute{
			"unraid": schema.StringAttribute{
				MarkdownDescription: "The Unraid OS version.",
				Computed:            true,
			},
			"api": schema.StringAttribute{
				MarkdownDescription: "The Unraid API version.",
				Computed:            true,
			},
			"kernel": schema.StringAttribute{
				MarkdownDescription: "The Linux kernel version.",
				Computed:            true,
			},
		},
	}
}

func (d *VersionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VersionDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	data, err := client.GetVersion(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Unraid version", err.Error())
		return
	}

	core := data.Info.Versions.Core
	state := VersionDataSourceModel{
		Unraid: types.StringValue(core.Unraid),
		Api:    types.StringValue(core.Api),
		Kernel: types.StringValue(core.Kernel),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
