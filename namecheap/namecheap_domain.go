package namecheap_provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

const (
	mode_create     = "create"
	mode_renew      = "renew"
	mode_reactivate = "reactivate"
	mode_skip       = "skip"
)

type namecheapDomainResource struct {
	client *namecheap.Client
}

type namecheapDomainState struct {
	Domain           types.String `tfsdk:"domain"`
	MinDaysRemaining types.Int64  `tfsdk:"min_days_remaining"`
	Years            types.Int64  `tfsdk:"auto_renew_years"`
}

func NewNamecheapDomainResource() resource.Resource {
	return &namecheapDomainResource{}
}

// Metadata
func (r *namecheapDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

// Schema
func (r *namecheapDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a domain in NameCheap",
		Attributes: map[string]schema.Attribute{
			"domain": &schema.StringAttribute{
				Description: "Domain name to manage in NameCheap",
				Required:    true,
			},
			"min_days_remaining": &schema.Int64Attribute{
				MarkdownDescription: "The minimum amount of days remaining on the expiration of a domain before a " +
					"renewal is attempted. The default is `30`. A value of less than `0` means that the domain will " +
					"never be renewed.",
				Optional: true,
				Default:  int64default.StaticInt64(30),
			},
			"auto_renew_years": &schema.Int64Attribute{
				MarkdownDescription: "Number of years to register and renew. The default is `1`. A value of less " +
					"than `0` means that the domain will never be auto renewed.",
				Optional: true,
				Default:  int64default.StaticInt64(1),
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *namecheapDomainResource) Configure(_ context.Context,
	req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// this data available on apply stage
		return
	}
	client, ok := req.ProviderData.(*namecheap.Client)
	if !ok {
		resp.Diagnostics.AddError("req.ProviderData not a namecheap.Client error", "")
		return
	}
	r.client = client
}

// Create
func (r *namecheapDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *namecheapDomainState
	d := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	mode := r.calculateMode(plan)
	log(ctx, "CalculateMode Complete,Mode = %s", mode)

	domain := plan.Domain.ValueString()
	years := plan.Years.ValueInt64()

	switch mode {
	case mode_create:
		diag := createDomain(ctx, domain, strconv.FormatInt(years, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case mode_renew:
		diag := renewDomain(ctx, domain, strconv.FormatInt(years, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case mode_reactivate:
		diag := reactivateDomain(ctx, domain, strconv.FormatInt(years, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case mode_skip:

	default:
		resp.Diagnostics.AddError("Invalid mode value", mode)
		return
	}

	state := &namecheapDomainState{
		Domain:           plan.Domain,
		Years:            plan.Years,
		MinDaysRemaining: plan.MinDaysRemaining,
	}
	setStateDiags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read
func (r *namecheapDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state *namecheapDomainState
	d := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(d...)
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
	var plan *namecheapDomainState
	// Retrieve values from plan
	d := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	newMode := r.calculateMode(plan)
	log(ctx, "CalculateMode Complete,Mode = %s", newMode)

	newDomain := plan.Domain.ValueString()
	newYear := plan.Years.ValueInt64()

	switch newMode {
	case mode_create:
		diag := createDomain(ctx, newDomain, strconv.FormatInt(newYear, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case mode_renew:
		diag := renewDomain(ctx, newDomain, strconv.FormatInt(newYear, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case mode_reactivate:
		diag := reactivateDomain(ctx, newDomain, strconv.FormatInt(newYear, 10), r.client)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case mode_skip:

	default:
		resp.Diagnostics.AddError("invalid mode value", newMode)
		return
	}

	// Set state items
	state := &namecheapDomainState{}
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

	var state *namecheapDomainState
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := state.Domain.ValueString()

	//since domain can not be deleted in NameCheap, so we do nothing here but give a warning
	msg := fmt.Sprintf("since domain can not be deleted in NameCheap, %s still exist actually  ", domain)
	tflog.Warn(ctx, msg)
}

func (r *namecheapDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import RecordId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}

func (r *namecheapDomainResource) calculateMode(plan *namecheapDomainState) string {
	domain := plan.Domain.ValueString()

	var req *namecheap.DomainsGetListArgs
	req.SearchTerm = &domain
	res, err := r.client.Domains.GetList(req)
	if err != nil {
		return mode_create
	}

	resName := *((*res.Domains)[0].Name)
	if resName != domain {
		return mode_create
	}

	isExpired := *((*res.Domains)[0].IsExpired)
	if isExpired {
		return mode_reactivate
	}

	minDaysRemain := plan.MinDaysRemaining.ValueInt64()
	expires := *((*res.Domains)[0].Expires)
	diff := time.Until(expires.Time)
	if int64(diff.Hours()) < minDaysRemain*24 {
		return mode_renew
	}

	return mode_skip
}
