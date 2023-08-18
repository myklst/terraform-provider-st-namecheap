package namecheap_provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

const MODE_CREATE = "create"
const MODE_RENEW = "renew"
const MODE_REACTIVATE = "reactivate"

func NewNamecheapDomainResource() resource.Resource {
	return &namecheapDomainResource{}
}

type namecheapDomainResource struct {
	client *namecheap.Client
}

type namecheapDomainResourceModel struct {
	Domain types.String `tfsdk:"domain"`
	Mode   types.String `tfsdk:"mode"`
	Years  types.Int64  `tfsdk:"years"`
}

// Metadata returns the resource namecheap_domain type name.
func (r *namecheapDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "namecheap_domain"
}

// Schema defines the schema for the namecheap_domain resource.
func (r *namecheapDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a namecheap_domain resource.",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Description: "Purchased available domain name on your account",
				Required:    true,
			},
			"mode": schema.StringAttribute{
				Description: "domain operation type, include create, renew, reactivate.",
				Required:    true,
			},
			"years": schema.Int64Attribute{
				Description: "Number of years to register",
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *namecheapDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*namecheap.Client)
}

/*
	func resourceNamecheapDomain() *schema.Resource {
		return &schema.Resource{
			CreateContext: resourceDomainCreate,
			UpdateContext: resourceDomainUpdate,
			ReadContext:   resourceDomainRead,
			DeleteContext: resourceDomainDelete,

			Importer: &schema.ResourceImporter{
				StateContext: resourceDomainImport,
			},

			Schema: map[string]*schema.Schema{
				"domain": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     false,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "Purchased available domain name on your account",
				},
				"mode": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "domain operation type, include create, renew, reactivate",
					DefaultFunc: schema.EnvDefaultFunc("NAMECHEAP_MODE", "CREATE"),
				},
				"years": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Number of years to register",
					Default:     "2",
				},
			},
		}
	}
*/

func (r *namecheapDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import RecordId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}

// Create a new namecheap_domain resource
func (r *namecheapDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	log(ctx, "[resourceDomainCreate!]")
	var plan *namecheapDomainResourceModel
	getStateDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := plan.Domain.ValueString()
	mode := plan.Mode.String()
	years := plan.Years.String()

	switch mode {
	case MODE_CREATE:
		diag := createDomain(ctx, domain, years, r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_RENEW:
		diag := renewDomain(ctx, domain, years, r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_REACTIVATE:
		diag := reactivateDomain(ctx, domain, years, r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	default:
		resp.Diagnostics.AddError("invalid mode value", mode)
		return
	}

	// Set state items
	state := &namecheapDomainResourceModel{}
	state.Mode = plan.Mode
	state.Domain = plan.Domain
	state.Years = plan.Years

	setStateDiags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read namecheap_domain resource information
func (r *namecheapDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	log(ctx, "[resourceDomainRead!]")
	// Get current state
	var state *namecheapDomainResourceModel
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()

	_, err := r.client.Domains.GetInfo(domain)

	if err == nil {
		setStateDiags := resp.State.Set(ctx, state)
		resp.Diagnostics.Append(setStateDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// Update namecheap_domain resource and sets the updated Terraform state on success.
func (r *namecheapDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	log(ctx, "[resourceRecordUpdate!]")
	var plan *namecheapDomainResourceModel
	// Retrieve values from plan
	getPlanDiags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(getPlanDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newDomain := plan.Domain.ValueString()
	newMode := plan.Mode.String()
	newYear := plan.Years.String()
	switch newMode {
	case MODE_CREATE:
		diag := createDomain(ctx, newDomain, newYear, r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_RENEW:
		diag := renewDomain(ctx, newDomain, newYear, r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_REACTIVATE:
		diag := reactivateDomain(ctx, newDomain, newYear, r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	default:
		resp.Diagnostics.AddError("invalid mode value", newMode)
		return
	}

}

// Delete namecheap_domain resource and removes the Terraform state on success.
func (r *namecheapDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	log(ctx, "[resourceRecordDelete!]")

	//since domain can not be deleted in namecheap, so we do nothing here
}
