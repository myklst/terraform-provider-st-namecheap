package namecheap_provider

import (
	"context"
	"errors"
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
			"auto_renew_years": &schema.Int64Attribute{
				MarkdownDescription: "Number of years to register and renew. The default is `1`. The value must greater than 0 and less than or equal to 10",
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

	domain := plan.Domain.ValueString()
	years := plan.Years.ValueInt64()

	diag := r.createDomain(ctx, domain, strconv.FormatInt(years, 10))
	resp.Diagnostics.Append(diag)
	if resp.Diagnostics.HasError() {
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
	} else {
		if strings.Contains(err.Error(), "Domain is invalid") {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Get domain info error ", err.Error())
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

	newMode, diag := r.calculateMode(ctx, plan)
	resp.Diagnostics.Append(diag)
	if resp.Diagnostics.HasError() {
		return
	}

	log(ctx, "CalculateMode Complete,Mode = %s", newMode)

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

func (r *namecheapDomainResource) calculateMode(ctx context.Context, plan *namecheapDomainState) (string, diag.Diagnostic) {
	domain := plan.Domain.ValueString()

	log(ctx, "calculateMode111")

	var req namecheap.DomainsGetListArgs

	req.SearchTerm = &domain

	res, err := r.client.Domains.GetList(&req)
	if err != nil {
		return "", DiagnosticErrorOf("domain [%s] doesn't exist", err, domain)
	}
	log(ctx, "calculateMode222")
	resName := *((*res.Domains)[0].Name)
	if resName != domain {
		return "", DiagnosticErrorOf("domain [%s] doesn't exist", nil, domain)
	}

	log(ctx, "calculateMode333")

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
	if int64(diff.Hours()) < minDaysRemain*24 {
		return mode_renew, nil
	}

	return mode_skip, nil
}

func DiagnosticErrorOf(format string, err error, a ...any) diag.Diagnostic {
	msg := fmt.Sprintf(format, a)
	if err != nil {
		return diag.NewErrorDiagnostic(msg, err.Error())
	} else {
		return diag.NewErrorDiagnostic(msg, "")
	}

}

func (r *namecheapDomainResource) renewDomain(ctx context.Context, domain string, years string) diag.Diagnostic {

	client := r.client
	resp, err := sdk.DomainsRenew(client, domain, years)

	if err != nil || *resp.Result.Renew == false {
		log(ctx, "renew domain %s failed, exit", domain)
		log(ctx, "reason:", err.Error())
		return DiagnosticErrorOf("renew domain [%s] failed", err, domain)
	}

	log(ctx, "renew domain [%s] success", domain)
	return nil
}

func (r *namecheapDomainResource) reactivateDomain(ctx context.Context, domain string, years string) diag.Diagnostic {

	client := r.client
	resp, err := sdk.DomainsReactivate(client, domain, years)

	if err != nil || !*resp.Result.IsSuccess {
		log(ctx, "reactivate domain %s failed, exit", domain)
		log(ctx, "reason:", err.Error())
		return DiagnosticErrorOf("reactivate domain [%s] failed", err, domain)
	}

	log(ctx, "reactivate domain [%s] success", domain)
	return nil
}

func (r *namecheapDomainResource) createDomain(ctx context.Context, domain string, years string) diag.Diagnostic {

	client := r.client
	//get domain info
	_, err := client.Domains.GetInfo(domain)
	if err == nil {
		return DiagnosticErrorOf("domain [%s] has been created in this account", nil, domain)
	}

	//if domain does not exist, then create

	//log.Println("Can not Get Domain Info, Creating:%s", domain)
	resp, err := sdk.DomainsAvailable(client, domain)
	if err == nil && *resp.Result.Available == true {
		// no err and available, create
		log(ctx, "Domain [%s] is available, Creating...", domain)
		r, err := r.getUserAccountContact()
		if err != nil {
			log(ctx, "get user Contacts failed, exit")
			log(ctx, "reason:", err.Error())
			return DiagnosticErrorOf("get user contacts failed, please check the contacts in the account manually", err, domain)
		}

		log(ctx, "debug domain create", r)

		_, err = sdk.DomainsCreate(client, domain, years, r)

		if err != nil {
			log(ctx, "create domain [%s] failed, exit", domain)
			log(ctx, "reason:", err.Error())
			return DiagnosticErrorOf("create domain [%s] failed", err, domain)
		}

	} else {
		log(ctx, "domain [%s] is not available, exiting!", domain)
		return DiagnosticErrorOf("domain [%s] is not available to register, you need to change to another domain", err, domain)
	}

	return nil

}

func (r *namecheapDomainResource) getUserAccountContact() (*sdk.UseraddrGetInfoCommandResponse, error) {

	client := r.client
	r1, err := sdk.UseraddrGetList(client)
	if err != nil {
		return nil, err
	}
	addrlen := len(*r1.Result.List)
	if addrlen == 0 {
		return nil, errors.New("UseraddrGetList returns 0, please add user contact info to this account")
	}

	addrId := *(*r1.Result.List)[0].AddressId

	r2, err := sdk.UseraddrGetInfo(client, addrId)
	if err != nil {
		return nil, err
	}

	return r2, nil

}

func log(ctx context.Context, format string, a ...any) {
	msg := fmt.Sprintf(format, a)
	tflog.Info(ctx, msg)
}
