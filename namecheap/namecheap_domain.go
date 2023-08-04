package namecheap_provider

import (
	"context"
	"strings"

	"github.com/agent-tao/go-namecheap-sdk/v2/namecheap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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
		},
	}
}

func resourceDomainImport(ctx context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log(ctx, "resourceRecordImport!!!!!!!!!!!!")
	if err := data.Set("domain", data.Id()); err != nil {

		return nil, err
	}

	return []*schema.ResourceData{data}, nil
}

func resourceDomainCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log(ctx, "resourceDomainCreate!!!!!!!!!!!!")
	client := meta.(*namecheap.Client)

	domain := strings.ToLower(data.Get("domain").(string))

	//create domain if Domain doesn't exist
	diags := createDomainIfNonexist(ctx, domain, client)
	if diags.HasError() {
		return diags
	}

	data.SetId(domain)

	return nil
}

func resourceDomainRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log(ctx, "resourceDomainRead!!!!!!!!!!!!")

	client := meta.(*namecheap.Client)

	domain := strings.ToLower(data.Get("domain").(string))

	_, err := client.Domains.GetInfo(domain)
	if err == nil {
		_ = data.Set("domain", domain)
	}

	return nil
}

func resourceDomainUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log(ctx, "resourceRecordUpdate!!!!!!!!!!!!")

	client := meta.(*namecheap.Client)

	//domain := strings.ToLower(data.Get("domain").(string))

	oldDomainRaw, newDomainRaw := data.GetChange("domain")

	oldDomain := oldDomainRaw.(string)
	newDomain := newDomainRaw.(string)

	if oldDomain != "" {
		//delete
	}

	if newDomain != "" {
		createDomainIfNonexist(ctx, newDomain, client)
	}

	return nil
}

func resourceDomainDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log(ctx, "resourceRecordDelete!!!!!!!!!!!!")

	// do nothing

	return nil
}
