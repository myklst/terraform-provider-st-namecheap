package namecheap_provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
	"os"
	"strconv"
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
				Description: "A registered user name for namecheap. May also be provided via NAMECHEAP_USER_NAME environment variable.",
				Optional:    true,
			},
			"api_user": schema.StringAttribute{
				Description: "A registered api user for namecheap. May also be provided via NAMECHEAP_API_USER environment variable.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "The namecheap API key. May also be provided via NAMECHEAP_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"client_ip": schema.StringAttribute{
				Description: "Client IP address. May also be provided via NAMECHEAP_CLIENT_IP environment variable.",
				Optional:    true,
			},
			"use_sandbox": schema.BoolAttribute{
				Description: "Use sandbox API endpoints. May also be provided via NAMECHEAP_USE_SANDBOX environment variable.",
				Optional:    true,
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

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.UserName.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("user_name"),
			"Unknown Namecheap user_name",
			"The provider cannot create the Namecheap API client as there is an unknown configuration value for the"+
				"Namecheap user_name. Set the value statically in the configuration, or use the NAMECHEAP_USER_NAME environment variable.",
		)
	}

	if config.ApiUser.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_user"),
			"Unknown Namecheap api_user",
			"The provider cannot create the Namecheap API client as there is an unknown configuration value for the"+
				"Namecheap api_user. Set the value statically in the configuration, or use the NAMECHEAP_API_USER environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Namecheap api_key",
			"The provider cannot create the Namecheap API client as there is an unknown configuration value for the"+
				"Namecheap api_key. Set the value statically in the configuration, or use the NAMECHEAP_API_KEY environment variable.",
		)
	}

	if config.ClientIp.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_ip"),
			"Unknown Namecheap client_ip",
			"The provider cannot create the Namecheap API client as there is an unknown configuration value for the"+
				"Namecheap client_ip. Set the value statically in the configuration, or use the NAMECHEAP_CLIENT_IP environment variable.",
		)
	}

	if config.UseSandbox.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("use_sandbox"),
			"Unknown Namecheap use_sandbox",
			"The provider cannot create the Namecheap API client as there is an unknown configuration value for the"+
				"Namecheap use_sandbox. Set the value statically in the configuration, or use the NAMECHEAP_USE_SANDBOX environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	var userName, apiUser, apiKey, clientIp string

	if !config.UserName.IsNull() {
		userName = config.UserName.ValueString()
	} else {
		userName = os.Getenv("NAMECHEAP_USER_NAME")
	}

	if !config.ApiUser.IsNull() {
		apiUser = config.ApiUser.ValueString()
	} else {
		apiUser = os.Getenv("NAMECHEAP_API_USER")
	}

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	} else {
		apiKey = os.Getenv("NAMECHEAP_API_KEY")
	}

	if !config.ClientIp.IsNull() {
		clientIp = config.ClientIp.ValueString()
	} else {
		clientIp = os.Getenv("NAMECHEAP_CLIENT_IP")
	}

	var useSandbox bool
	if !config.UseSandbox.IsNull() {
		useSandbox = config.UseSandbox.ValueBool()
	} else {
		var err error
		useSandbox, err = strconv.ParseBool(os.Getenv("NAMECHEAP_USE_SANDBOX"))
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("use_sandbox"),
				"Missing Namecheap use_sandbox",
				"The provider cannot create the Namecheap API client as there is a "+
					"missing or empty value for the Namecheap API use_sandbox. Set the "+
					"use_sandbox value in the configuration or use the NAMECHEAP_USE_SANDBOX "+
					"environment variable. If either is already set, ensure the value "+
					"is not empty.",
			)
		}
	}

	// If any of the expected configuration are missing, return
	// errors with provider-specific guidance.
	if userName == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("user_name"),
			"Missing Namecheap API user_name",
			"The provider cannot create the Namecheap API client as there is a "+
				"missing or empty value for the Namecheap API user_name. Set the "+
				"user_name value in the configuration or use the NAMECHEAP_USER_NAME "+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if apiUser == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_user"),
			"Missing Namecheap API access key",
			"The provider cannot create the Namecheap API client as there is a "+
				"missing or empty value for the Namecheap API api_user. Set the "+
				"api_user value in the configuration or use the NAMECHEAP_API_USER "+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Namecheap secret key",
			"The provider cannot create the Namecheap API client as there is a "+
				"missing or empty value for the Namecheap API api_key. Set the "+
				"api_key value in the configuration or use the NAMECHEAP_API_KEY "+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if clientIp == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_ip"),
			"Missing Namecheap secret key",
			"The provider cannot create the Namecheap API client as there is a "+
				"missing or empty value for the Namecheap API client_ip. Set the "+
				"client_ip value in the configuration or use the NAMECHEAP_CLIENT_IP "+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

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
