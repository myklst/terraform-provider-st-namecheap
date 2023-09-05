package namecheap_provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

const MODE_CREATE = "create"
const MODE_RENEW = "renew"
const MODE_REACTIVATE = "reactivate"
const MODE_SKIP = "skip"

func NewNamecheapDomainResource() resource.Resource {
	return &namecheapDomainResource{}
}

type namecheapDomainResource struct {
	client *namecheap.Client
}

type namecheapDomainResourceModel struct {
	Domain           types.String `tfsdk:"domain"`
	MinDaysRemaining types.Int64  `tfsdk:"min_days_remaining"`
	Years            types.Int64  `tfsdk:"auto_renew_years"`
}

// Metadata returns the resource namecheap_domain type name.
func (r *namecheapDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

// Schema defines the schema for the namecheap_domain resource.
func (r *namecheapDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Description: "Domain name to create",
				Required:    true,
			},
			"min_days_remaining": schema.Int64Attribute{
				Description: "maintain min days remain to renew, reactivate.",
				Required:    true,
			},
			"auto_renew_years": schema.Int64Attribute{
				Description: "Number of years to register and renew",
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
	years := plan.Years.ValueInt64()

	mode := CalculateMode(r.client, plan)
	log(ctx, "CalculateMode Complete,Mode = %s", mode)

	switch mode {
	case MODE_CREATE:
		diag := createDomain(ctx, domain, strconv.FormatInt(years, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_RENEW:
		diag := renewDomain(ctx, domain, strconv.FormatInt(years, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_REACTIVATE:
		diag := reactivateDomain(ctx, domain, strconv.FormatInt(years, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_SKIP:

	default:
		resp.Diagnostics.AddError("invalid mode value", mode)
		return
	}

	// Set state items
	state := &namecheapDomainResourceModel{}
	state.Domain = plan.Domain
	state.Years = plan.Years
	state.MinDaysRemaining = plan.MinDaysRemaining

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

	newYear := plan.Years.ValueInt64()
	newMode := CalculateMode(r.client, plan)
	log(ctx, "CalculateMode Complete,Mode = %s", newMode)

	switch newMode {
	case MODE_CREATE:
		diag := createDomain(ctx, newDomain, strconv.FormatInt(newYear, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_RENEW:
		diag := renewDomain(ctx, newDomain, strconv.FormatInt(newYear, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_REACTIVATE:
		diag := reactivateDomain(ctx, newDomain, strconv.FormatInt(newYear, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_SKIP:

	default:
		resp.Diagnostics.AddError("invalid mode value", newMode)
		return
	}

	// Set state items
	state := &namecheapDomainResourceModel{}
	state.Domain = plan.Domain
	state.Years = plan.Years
	state.MinDaysRemaining = plan.MinDaysRemaining

	setStateDiags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete namecheap_domain resource and removes the Terraform state on success.
func (r *namecheapDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	log(ctx, "[resourceRecordDelete!]")

	var state *namecheapDomainResourceModel
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := state.Domain.ValueString()

	//since domain can not be deleted in namecheap, so we do nothing here but give a warning
	msg := fmt.Sprintf("since domain can not be deleted in namecheap, %s still exist actually  ", domain)
	tflog.Warn(ctx, msg)
}
