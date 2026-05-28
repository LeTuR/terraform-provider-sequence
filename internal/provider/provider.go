// Copyright (c) Arthur Cesaré-Herriau
// SPDX-License-Identifier: MIT

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &SequenceProvider{}

type SequenceProvider struct {
	version string
}

func (p *SequenceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sequence"
	resp.Version = p.version
}

func (p *SequenceProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The sequence provider generates zero-padded sequential numbers gated by a trigger.",
	}
}

func (p *SequenceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *SequenceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNumberResource,
	}
}

func (p *SequenceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SequenceProvider{version: version}
	}
}
