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

type infoQueryResponse struct {
	Info struct {
		Versions struct {
			Core struct {
				Unraid string `json:"unraid"`
				Api    string `json:"api"`
				Kernel string `json:"kernel"`
			} `json:"core"`
		} `json:"versions"`
	} `json:"info"`
}

func (d *VersionDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	const query = `
		query {
			info {
				versions {
					core {
						unraid
						api
						kernel
					}
				}
			}
		}
	`

	data, err := client.Do[infoQueryResponse](ctx, d.client, query, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Unraid info", err.Error())
		return
	}

	state := VersionDataSourceModel{
		Unraid: types.StringValue(data.Info.Versions.Core.Unraid),
		Api:    types.StringValue(data.Info.Versions.Core.Api),
		Kernel: types.StringValue(data.Info.Versions.Core.Kernel),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
