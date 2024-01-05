package namecheap

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"

	"github.com/myklst/terraform-provider-st-namecheap/namecheap/sdk"
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
	Years            types.Int64  `tfsdk:"purchase_years"`
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"min_days_remaining": &schema.Int64Attribute{
				MarkdownDescription: "The minimum amount of days remaining on the expiration of a domain before a " +
					"renewal is attempted. The default is `30`. A value of less than `0` means that the domain will " +
					"never be renewed.",
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(30),
			},
			"purchase_years": &schema.Int64Attribute{
				MarkdownDescription: "Number of years to purchase and renew. The default is `1`. The value must greater than 0 and less than or equal to 10",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
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
		resp.Diagnostics.AddError("req.ProviderData isn't a namecheap.Client", "")
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

	domain := plan.Domain.ValueString()
	years := plan.Years.ValueInt64()

	d1 := r.createDomain(ctx, domain, strconv.FormatInt(years, 10))
	resp.Diagnostics.Append(d1)
	if resp.Diagnostics.HasError() {
		return
	}

	state := namecheapDomainState{
		Domain:           plan.Domain,
		Years:            plan.Years,
		MinDaysRemaining: plan.MinDaysRemaining,
	}
	d2 := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(d2...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read
func (r *namecheapDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *namecheapDomainState
	d := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()
	if _, err := r.client.Domains.GetInfo(domain); err != nil {
		if strings.Contains(err.Error(), "Domain is invalid") {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Get domain info error ", err.Error())
		}
		return
	}

	d1 := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(d1...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update namecheap_domain resource and sets the updated Terraform state on success.
func (r *namecheapDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *namecheapDomainState
	d := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	newMode, diag := r.calculateMode(ctx, plan)
	resp.Diagnostics.Append(diag)
	if resp.Diagnostics.HasError() {
		return
	}

	log(ctx, "CalculateMode result = %s", newMode)

	newDomain := plan.Domain.ValueString()
	newYear := plan.Years.ValueInt64()

	switch newMode {
	case mode_renew:
		diag := r.renewDomain(ctx, newDomain, strconv.FormatInt(newYear, 10))
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case mode_reactivate:
		diag := r.reactivateDomain(ctx, newDomain, strconv.FormatInt(newYear, 10))
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	case mode_skip:

	default:
		resp.Diagnostics.AddError("invalid mode value", newMode)
		return
	}

	// Set state
	state := namecheapDomainState{
		Domain:           plan.Domain,
		Years:            plan.Years,
		MinDaysRemaining: plan.MinDaysRemaining,
	}
	setStateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete namecheap_domain resource and removes the Terraform state on success.
func (r *namecheapDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *namecheapDomainState
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()
	// Since domain can not be deleted in NameCheap, so we do nothing here but give a warning
	msg := fmt.Sprintf("Since domain can not be deleted in NameCheap, %s still exist actually", domain)
	tflog.Warn(ctx, msg)
}

func (r *namecheapDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import RecordId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}

func (r *namecheapDomainResource) calculateMode(ctx context.Context, plan *namecheapDomainState) (string, diag.Diagnostic) {
	domain := plan.Domain.ValueString()
	res, err := r.client.Domains.GetList(&namecheap.DomainsGetListArgs{
		SearchTerm: &domain,
	})
	if err != nil {
		return "", diagnosticErrorOf(err, "domain [%s] doesn't exist", domain)
	}

	respName := *((*res.Domains)[0].Name)
	if respName != domain {
		return "", diagnosticErrorOf(nil, "domain [%s] doesn't exist", domain)
	}

	minDaysRemain := plan.MinDaysRemaining.ValueInt64()
	if minDaysRemain <= 0 {
		return mode_skip, nil
	}

	isExpired := *((*res.Domains)[0].IsExpired)
	if isExpired {
		return mode_reactivate, nil
	}

	expires := *((*res.Domains)[0].Expires)
	diff := time.Until(expires.Time)
	if int64(diff.Hours())/24 < minDaysRemain {
		return mode_renew, nil
	}

	return mode_skip, nil
}

func (r *namecheapDomainResource) createDomain(ctx context.Context, domain string, years string) diag.Diagnostic {
	client := r.client
	// Get domain info
	if _, err := client.Domains.GetInfo(domain); err == nil {
		return diagnosticErrorOf(nil, "domain [%s] has been created in this account", domain)
	}
	// else, if domain does not exist, then create

	resp, err := sdk.DomainsAvailable(client, domain)
	if err == nil && resp.Result.Available {
		// no err and available, create
		log(ctx, "Domain [%s] is available, Creating...", domain)

		r, err := r.getUserAccountContact()
		if err != nil {
			log(ctx, "get user contacts failed: %s", err.Error())
			return diagnosticErrorOf(err, "get user contacts failed: %s", domain)
		}

		_, err = sdk.DomainsCreate(client, domain, years, r)
		if err != nil {
			log(ctx, "create domain [%s] failed: %s", domain, err.Error())
			return diagnosticErrorOf(err, "create domain [%s] failed", domain)
		}
	} else {
		log(ctx, "domain [%s] is not available, exiting!", domain)
		return diagnosticErrorOf(err, "domain [%s] is not available to register, you need to change to another domain", domain)
	}

	return nil
}

func (r *namecheapDomainResource) renewDomain(ctx context.Context, domain string, years string) diag.Diagnostic {
	client := r.client
	resp, err := sdk.DomainsRenew(client, domain, years)

	if err != nil || !resp.Result.Renew {
		log(ctx, "renew domain %s failed, exit", domain)
		log(ctx, "reason: %s", err.Error())
		return diagnosticErrorOf(err, "renew domain [%s] failed", domain)
	}

	log(ctx, "renew domain [%s] success", domain)
	return nil
}

func (r *namecheapDomainResource) reactivateDomain(ctx context.Context, domain string, years string) diag.Diagnostic {
	client := r.client
	resp, err := sdk.DomainsReactivate(client, domain, years)

	if err != nil || !resp.Result.IsSuccess {
		log(ctx, "reactivate domain %s failed: %s", domain, err.Error())
		return diagnosticErrorOf(err, "reactivate domain [%s] failed", domain)
	}

	log(ctx, "reactivate domain [%s] success", domain)
	return nil
}

func (r *namecheapDomainResource) getUserAccountContact() (*sdk.UserAddrGetInfoCommandResponse, error) {
	client := r.client

	// r1, err := sdk.UserAddrGetList(client)
	// if err != nil {
	// 	return nil, err
	// }
	// addrlen := len(*r1.Result.List)
	// if addrlen == 0 {
	// 	return nil, errors.New("UseraddrGetList returns 0, please add user contact info to this account")
	// }
	// addrId := *(*r1.Result.List)[0].AddressId

	// Hard-coded addrId to "0", to use `Primary Address`
	r2, err := sdk.UserAddrGetInfo(client, "0")
	if err != nil {
		return nil, err
	}

	return r2, nil
}

func log(ctx context.Context, format string, a ...any) {
	tflog.Info(ctx, fmt.Sprintf(format, a...))
}

func diagnosticErrorOf(err error, format string, a ...any) diag.Diagnostic {
	msg := fmt.Sprintf(format, a...)
	if err != nil {
		return diag.NewErrorDiagnostic(msg, err.Error())
	} else {
		return diag.NewErrorDiagnostic(msg, "")
	}
}
