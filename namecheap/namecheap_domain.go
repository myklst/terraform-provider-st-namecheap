package namecheap_provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

const MODE_CREATE = "create"
const MODE_RENEW = "renew"
const MODE_REACTIVATE = "reactivate"

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

func resourceDomainImport(ctx context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log(ctx, "[resourceRecordImport!]")
	if err := data.Set("domain", data.Id()); err != nil {
		return nil, err
	}
	if err := data.Set("mode", MODE_CREATE); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{data}, nil
}

func resourceDomainCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log(ctx, "[resourceDomainCreate!]")
	client := meta.(*namecheap.Client)

	domain := strings.ToLower(data.Get("domain").(string))
	mode := strings.ToLower(data.Get("mode").(string))
	years := data.Get("years").(string)

	switch mode {
	case MODE_CREATE:
		diags := createDomain(ctx, domain, years, client)
		if diags.HasError() {
			return diags
		}
	case MODE_RENEW:
		diags := renewDomain(ctx, domain, years, client)
		if diags.HasError() {
			return diags
		}
	case MODE_REACTIVATE:
		diags := reactivateDomain(ctx, domain, years, client)
		if diags.HasError() {
			return diags
		}
	default:
		//do nothing
	}

	data.SetId(domain)

	return nil
}

func resourceDomainRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log(ctx, "[resourceDomainRead!]")

	client := meta.(*namecheap.Client)

	domain := strings.ToLower(data.Get("domain").(string))

	_, err := client.Domains.GetInfo(domain)
	if err == nil {
		_ = data.Set("domain", domain)
	}

	return nil
}

func resourceDomainUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log(ctx, "[resourceRecordUpdate!]")

	client := meta.(*namecheap.Client)

	//we can do nothing on old name,year and mode
	_, newDomainRaw := data.GetChange("domain")
	newDomain := newDomainRaw.(string)

	_, newYearRaw := data.GetChange("years")
	newYear := newYearRaw.(string)

	_, newModeRaw := data.GetChange("mode")
	newMode := newModeRaw.(string)

	switch newMode {
	case MODE_CREATE:
		diags := createDomain(ctx, newDomain, newYear, client)
		if diags.HasError() {
			return diags
		}
	case MODE_RENEW:
		diags := renewDomain(ctx, newDomain, newYear, client)
		if diags.HasError() {
			return diags
		}
	case MODE_REACTIVATE:
		diags := reactivateDomain(ctx, newDomain, newYear, client)
		if diags.HasError() {
			return diags
		}
	default:

	}

	return nil
}

func resourceDomainDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log(ctx, "[resourceRecordDelete!]")

	// do nothing

	return nil
}
