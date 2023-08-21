package namecheap_provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type namecheapProviderModel struct {
	UserName   types.String `tfsdk:"user_name"`
	ApiUser    types.String `tfsdk:"api_user"`
	ApiKey     types.String `tfsdk:"api_key"`
	ClientIp   types.String `tfsdk:"client_ip"`
	UseSandbox types.Bool   `tfsdk:"use_sandbox"`
}

// Metadata returns the provider type name.
func (p *namecheapProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "st-namecheap"
}

// Schema defines the provider-level schema for configuration data.
func (p *namecheapProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The namecheap domain provider is used to interact with the namecheap to manage domains from it. " +
			"The provider needs to be configured with the proper credentials before it can be used.",
		Attributes: map[string]schema.Attribute{
			"user_name": schema.StringAttribute{
				Description: "A registered user name for namecheap",
				Required:    true,
			},
			"api_user": schema.StringAttribute{
				Description: "A registered api user for namecheap",
				Required:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "The namecheap API key",
				Required:    true,
				Sensitive:   true,
			},
			"client_ip": schema.StringAttribute{
				Description: "Client IP address",
				Required:    true,
			},
			"use_sandbox": schema.BoolAttribute{
				Description: "Use sandbox API endpoints",
				Required:    true,
			},
		},
	}
}

func (p *namecheapProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var config namecheapProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userName := config.UserName.ValueString()
	apiUser := config.ApiUser.ValueString()
	apiKey := config.ApiKey.ValueString()
	clientIp := config.ClientIp.ValueString()
	useSandbox := config.UseSandbox.ValueBool()

	client := namecheap.NewClient(&namecheap.ClientOptions{
		UserName:   userName,
		ApiUser:    apiUser,
		ApiKey:     apiKey,
		ClientIp:   clientIp,
		UseSandbox: useSandbox,
	})

	resp.ResourceData = client
}

func (p *namecheapProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNamecheapDomainResource,
	}
}

func (p *namecheapProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &namecheapProvider{}
)

// New is a helper function to simplify provider server
func New() provider.Provider {
	return &namecheapProvider{}
}

type namecheapProvider struct{}
