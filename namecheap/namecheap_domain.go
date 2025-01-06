package namecheap

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"

	"github.com/myklst/terraform-provider-st-namecheap/namecheap/sdk"
)

const (
	MODE_RENEW      string = "renew"
	MODE_REACTIVATE string = "reactivate"
)

type namecheapDomainResource struct {
	client *namecheap.Client
}

type namecheapDomainState struct {
	Domain           types.String  `tfsdk:"domain"`
	Nameservers      types.List    `tfsdk:"nameservers"`
	MaxPrice         types.Float64 `tfsdk:"max_price"`
	MinDaysRemaining types.Int64   `tfsdk:"min_days_remaining"`
	Years            types.Int64   `tfsdk:"purchase_years"`
	DomainExpiryDate types.String  `tfsdk:"domain_expiry_date"`
	RequiredRenew    types.Bool    `tfsdk:"required_renew"`
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
				MarkdownDescription: "Domain name to manage in NameCheap",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nameservers": &schema.ListAttribute{
				MarkdownDescription: "Nameservers for the domain",
				Required:            true,
				ElementType:         types.StringType,
			},
			"max_price": &schema.Float64Attribute{
				MarkdownDescription: "Maximum price of the purchase domain",
				Required:            true,
			},
			"min_days_remaining": &schema.Int64Attribute{
				MarkdownDescription: "The minimum amount of days remaining on the expiration of a domain before a " +
					"renewal is attempted. The default is `30`. A value of less than `0` means that the domain will " +
					"never be renewed.",
				Optional: true,
			},
			"purchase_years": &schema.Int64Attribute{
				MarkdownDescription: "Number of years to purchase and renew. The default is `1`. The value must greater than 0 and less than or equal to 10",
				Optional:            true,
			},
			"domain_expiry_date": &schema.StringAttribute{
				MarkdownDescription: "The expiry date of the domain, stored in ISO 8601 format (e.g., `2024-12-30T14:59:59Z`). This field is computed automatically based on the domain's expiration date.",
				Computed:            true,
			},
			"required_renew": &schema.BoolAttribute{
				MarkdownDescription: "A boolean flag to keep track of whether domain renewal action is required. ",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *namecheapDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	maxprice := plan.MaxPrice.ValueFloat64()
	var nameservers string
	for _, x := range plan.Nameservers.Elements() {
		nameservers += strings.Trim(x.String(), "\"") + ","
	}

	createDomain := func() error {
		d1 := r.createDomain(ctx, domain, strconv.FormatInt(years, 10), nameservers, maxprice)
		resp.Diagnostics.Append(d1)
		if resp.Diagnostics.HasError() {
			return fmt.Errorf("domain creation failed: %v", resp.Diagnostics)
		}
		return nil
	}

	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = 30 * time.Second
	err := backoff.Retry(createDomain, reconnectBackoff)
	if err != nil {
		resp.Diagnostics.Append(diagnosticErrorOf(err, "domain [%s] creation failed after retries", domain))
		return
	}

	state := namecheapDomainState{
		Domain:           plan.Domain,
		Years:            plan.Years,
		MaxPrice:         plan.MaxPrice,
		MinDaysRemaining: plan.MinDaysRemaining,
		Nameservers:      plan.Nameservers,
	}

	// Compute `domainExpiryDate` and `domainExpiryRemainingDays` to get the expiration date and
	// remaining active days of the domain.
	domainExpiryDate, _err := r.getDomainExpiryDate(plan.Domain.ValueString())
	if _err != nil {
		resp.Diagnostics.Append(_err)
		return
	}
	state.DomainExpiryDate = types.StringValue(domainExpiryDate.Format("2006-01-02T15:04:05Z"))
	state.RequiredRenew = types.BoolValue(false)

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
	getResp, err := r.client.Domains.GetInfo(domain)
	if err != nil {
		if strings.Contains(err.Error(), "Domain is invalid") {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Get domain info error ", err.Error())
		}
		return
	}

	nameserver := []attr.Value{}
	for _, x := range *getResp.DomainDNSGetListResult.DnsDetails.Nameservers {
		nameserver = append(nameserver, types.StringValue(x))
	}
	state.Nameservers = types.ListValueMust(types.StringType, nameserver)

	domainExpiryDate, _err := r.getDomainExpiryDate(domain)
	if _err != nil {
		resp.Diagnostics.Append(_err)
		return
	}
	state.DomainExpiryDate = types.StringValue(domainExpiryDate.Format("2006-01-02T15:04:05Z"))

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

	// Set state
	state := namecheapDomainState{
		Domain:           plan.Domain,
		Years:            plan.Years,
		MaxPrice:         plan.MaxPrice,
		MinDaysRemaining: plan.MinDaysRemaining,
		Nameservers:      plan.Nameservers,
	}

	// Compute `domainExpiryDate` and `domainExpiryRemainingDays` to get the expiration date and
	// remaining active days of the domain.
	domainExpiryDate, err := r.getDomainExpiryDate(plan.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.Append(err)
		return
	}

	// Calculate domain active remaining days
	domainExpiryRemainingDays := calcDomainRemainingDays(domainExpiryDate)

	// Attempt to renew / reactivate domain if the `DomainRemainingDays` is lesser or equal to `MinDaysRemaining`
	if domainExpiryRemainingDays <= plan.MinDaysRemaining.ValueInt64() {
		domain := plan.Domain.ValueString()
		renewYear := plan.Years.ValueInt64()

		newMode, diag := r.calculateMode(domain)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}

		switch newMode {
		case MODE_RENEW:
			diag := r.renewDomain(ctx, domain, strconv.FormatInt(renewYear, 10))
			resp.Diagnostics.Append(diag)
			if resp.Diagnostics.HasError() {
				return
			}
		case MODE_REACTIVATE:
			diag := r.reactivateDomain(ctx, domain, strconv.FormatInt(renewYear, 10))
			resp.Diagnostics.Append(diag)
			if resp.Diagnostics.HasError() {
				return
			}
		default:
			resp.Diagnostics.AddError("Invalid mode value", newMode)
			return
		}

		// Update and refresh expiration details after domain renewal / reactivate is done.
		domainExpiryDate, err = r.getDomainExpiryDate(plan.Domain.ValueString())
		if err != nil {
			resp.Diagnostics.Append(err)
			return
		}
	}

	// Update and refresh state for attributes `domainExpiryDate` & `domainExpirationDays`
	state.DomainExpiryDate = types.StringValue(domainExpiryDate.Format("2006-01-02T15:04:05Z"))
	state.RequiredRenew = types.BoolValue(false)

	// Configure nameservers
	var nameservers []string
	for _, x := range plan.Nameservers.Elements() {
		nameservers = append(nameservers, strings.Trim(x.String(), "\""))
	}
	_, _err := r.client.DomainsDNS.SetCustom(plan.Domain.ValueString(), nameservers)
	if _err != nil {
		resp.Diagnostics.AddError("Set nameserver failed error ", _err.Error())
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
	msg := fmt.Sprintf("Since domain can not be deleted in NameCheap, [%s] still exist actually", domain)
	tflog.Warn(ctx, msg)
}

func (r *namecheapDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import RecordId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}

func (r *namecheapDomainResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !req.Plan.Raw.IsNull() {
		var plan *namecheapDomainState
		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan namecheapDomainState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	requiresRenew, err := isDomainRequiredRenew(plan.MinDaysRemaining.ValueInt64(), plan.DomainExpiryDate.ValueString())
	if err != nil {
		return
	}

	// Enforce update lifecycle if domain is required to renew
	if requiresRenew {
		plan.DomainExpiryDate = types.StringUnknown()
		plan.RequiredRenew = types.BoolUnknown()

		setDomainExpiryDate := resp.Plan.SetAttribute(ctx, path.Root("domain_expiry_date"), types.StringUnknown())
		setRequiredRenew := resp.Plan.SetAttribute(ctx, path.Root("required_renew"), types.BoolUnknown())

		resp.Diagnostics.Append(setDomainExpiryDate...)
		resp.Diagnostics.Append(setRequiredRenew...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Plan.Set(ctx, plan)
}

func (r *namecheapDomainResource) calculateMode(domain string) (string, diag.Diagnostic) {
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

	isExpired := *((*res.Domains)[0].IsExpired)
	if isExpired {
		return MODE_REACTIVATE, nil
	}

	return MODE_RENEW, nil
}

func (r *namecheapDomainResource) createDomain(ctx context.Context, domain string, years string, nameservers string, maxprice float64) diag.Diagnostic {
	client := r.client
	// Get domain info
	if _, err := client.Domains.GetInfo(domain); err == nil {
		return diagnosticErrorOf(nil, "domain [%s] has been created in this account", domain)
	}

	// else, if domain does not exist, check for pricing then create
	resp, err := sdk.DomainsAvailable(client, domain)
	if err == nil && resp.Result.Available {
		// Check domain price before proceed
		var price float64

		// Check if the domain is a premium domain
		if resp.Result.Price != "0" {
			price, err = strconv.ParseFloat(resp.Result.Price, 32)
			if err != nil {
				return diagnosticErrorOf(err, "get domain price failed: %s", domain)
			}
		} else { //Do a normal price query on the target TLD
			priceResp, err := sdk.UserGetPricing(client, "register", domain)
			if err == nil {
				for _, s := range priceResp.Result.ProductCategory.Price {
					if s.Duration == years {
						if price, err = strconv.ParseFloat(s.Price, 32); err != nil {
							return diagnosticErrorOf(err, "get domain price failed: %s", domain)
						}
					}
				}
			} else {
				return diagnosticErrorOf(err, "get domain price failed: %s", domain)
			}
		}

		if price <= maxprice {
			// no err, price ok and available, create
			log(ctx, "Domain [%s] is available, Creating...", domain)

			r, err := r.getUserAccountContact()
			if err != nil {
				log(ctx, "get user contacts failed: %s", err.Error())
				return diagnosticErrorOf(err, "get user contacts failed: %s", domain)
			}

			_, err = sdk.DomainsCreate(client, domain, years, nameservers, r)
			if err != nil {
				log(ctx, "create domain [%s] failed: %s", domain, err.Error())
				return diagnosticErrorOf(err, "create domain [%s] failed", domain)
			}
		} else {
			log(ctx, "domain [%s] is overprice, exiting!", domain)
			return diagnosticErrorOf(err, "domain [%s] is overprice [%f], you need to change to another domain", domain, price)
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

func (r *namecheapDomainResource) getDomainExpiryDate(domain string) (time.Time, diag.Diagnostic) {
	var domainExpiryDate time.Time
	var err error

	getDomainExpiryInfo := func() error {
		res, err := r.client.Domains.GetList(&namecheap.DomainsGetListArgs{
			SearchTerm: &domain,
		})
		if err != nil {
			return fmt.Errorf("domain [%s] doesn't exist: %v", domain, err)
		}

		domainExpiryDate = (*res.Domains)[0].Expires.Time

		return nil
	}

	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = 30 * time.Second
	_err := backoff.Retry(getDomainExpiryInfo, reconnectBackoff)

	if _err != nil {
		return time.Time{}, diagnosticErrorOf(err, "failed to fetch domain expiry for [%s] after retries", domain)
	}

	return domainExpiryDate, nil
}

// Function to calculate the remaining active days of domain
func calcDomainRemainingDays(expiryDate time.Time) int64 {
	todayDate := time.Now()

	return int64(expiryDate.Sub(todayDate).Hours() / 24)
}

// Function to determine when the domain requires renewal action
func isDomainRequiredRenew(minDaysremaining int64, domainExpiry string) (bool, error) {
	domainExpiryDate, err := time.Parse("2006-01-02T15:04:05Z", domainExpiry)
	if err != nil {
		return false, err
	}

	domainRemainingDays := calcDomainRemainingDays(domainExpiryDate)

	return (domainRemainingDays <= minDaysremaining), err
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
