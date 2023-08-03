package namecheap_provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/agent-tao/go-namecheap-sdk/v2/namecheap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNamecheapDomainRecords() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordCreate,
		UpdateContext: resourceRecordUpdate,
		ReadContext:   resourceRecordRead,
		DeleteContext: resourceRecordDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Purchased available domain name on your account",
			},
			"email_type": {
				ConflictsWith: []string{"nameservers"},
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringInSlice(namecheap.AllowedEmailTypeValues, false),
				Description:   fmt.Sprintf("Possible values: %s", strings.TrimSpace(strings.Join(namecheap.AllowedEmailTypeValues, ", "))),
			},
			"record": {
				ConflictsWith: []string{"nameservers"},
				Type:          schema.TypeSet,
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Sub-domain/hostname to create the record for",
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(namecheap.AllowedRecordTypeValues, false),
							Description:  fmt.Sprintf("Possible values: %s", strings.TrimSpace(strings.Join(namecheap.AllowedRecordTypeValues, ", "))),
						},
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Possible values are URL or IP address. The value for this parameter is based on record type",
						},
						"mx_pref": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     10,
							Description: "MX preference for host. Applicable for MX records only",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1799,
							Description: fmt.Sprintf("Time to live for all record types. Possible values: any value between %d to %d", namecheap.MinTTL, namecheap.MaxTTL),
						},
					},
				},
			},
			"nameservers": {
				ConflictsWith: []string{"email_type", "record"},
				Type:          schema.TypeSet,
				Optional:      true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceRecordImport(ctx context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log(ctx, "resourceRecordImport!!!!!!!!!!!!")
	if err := data.Set("domain", data.Id()); err != nil {

		return nil, err
	}

	return []*schema.ResourceData{data}, nil
}

func resourceRecordCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log(ctx, "resourceRecordCreate!!!!!!!!!!!!")
	client := meta.(*namecheap.Client)

	domain := strings.ToLower(data.Get("domain").(string))

	//create domain if Domain doesn't exist
	diags := createDomainIfNonexist(ctx, domain, client)
	if diags.HasError() {
		return diags
	}

	var emailType *string
	var records []interface{}
	var nameservers []interface{}

	if emailTypeRaw, ok := data.GetOk("email_type"); ok {
		emailTypeString := emailTypeRaw.(string)
		emailType = &emailTypeString
	}

	if recordsRaw, ok := data.GetOk("record"); ok {
		records = recordsRaw.(*schema.Set).List()
	}

	if nameserversRaw, ok := data.GetOk("nameservers"); ok {
		nameservers = nameserversRaw.(*schema.Set).List()
	}

	if records != nil {
		diags := createRecordsOverwrite(ctx, domain, emailType, records, client)
		if diags.HasError() {
			return diags
		}
	}

	if nameservers != nil {
		diags := createNameserversOverwrite(ctx, domain, convertInterfacesToString(nameservers), client)
		if diags.HasError() {
			return diags
		}
	}

	data.SetId(domain)

	return nil
}

func resourceRecordRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log(ctx, "resourceRecordRead!!!!!!!!!!!!")

	client := meta.(*namecheap.Client)

	domain := strings.ToLower(data.Get("domain").(string))

	var emailType *string
	var records []interface{}
	var nameservers []interface{}

	if emailTypeRaw, ok := data.GetOk("email_type"); ok {
		emailTypeString := emailTypeRaw.(string)
		emailType = &emailTypeString
	}

	if recordsRaw, ok := data.GetOk("record"); ok {
		records = recordsRaw.(*schema.Set).List()
	}

	if nameserversRaw, ok := data.GetOk("nameservers"); ok {
		nameservers = nameserversRaw.(*schema.Set).List()
	}
	// We must read nameservers status before hosts.
	// If you're using custom nameservers, then the reading records process will fail since Namecheap doesn't control
	// the domain behaviour.
	nsResponse, err := client.DomainsDNS.GetList(domain)
	if err != nil {
		return diag.FromErr(err)
	}

	if !*nsResponse.DomainDNSGetListResult.IsUsingOurDNS {

		realNameservers, diags := readNameserversOverwrite(ctx, domain, client)
		if diags.HasError() {
			return diags
		}
		_ = data.Set("nameservers", *realNameservers)

		_ = data.Set("record", []interface{}{})
	} else {

		realRecords, realEmailType, diags := readRecordsOverwrite(domain, records, client)

		if diags.HasError() {
			return diags
		}

		_ = data.Set("record", *realRecords)
		if emailType != nil {
			_ = data.Set("email_type", *realEmailType)
		}

		if nameservers != nil {
			_ = data.Set("nameservers", []string{})
		}

	}

	return nil
}

func resourceRecordUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log(ctx, "resourceRecordUpdate!!!!!!!!!!!!")

	client := meta.(*namecheap.Client)

	domain := strings.ToLower(data.Get("domain").(string))

	oldRecordsRaw, newRecordsRaw := data.GetChange("record")
	oldNameserversRaw, newNameserversRaw := data.GetChange("nameservers")

	oldRecords := oldRecordsRaw.(*schema.Set).List()
	newRecords := newRecordsRaw.(*schema.Set).List()

	oldNameservers := oldNameserversRaw.(*schema.Set).List()
	newNameservers := newNameserversRaw.(*schema.Set).List()

	oldRecordsLen := len(oldRecords)
	newRecordsLen := len(newRecords)

	oldNameserversLen := len(oldNameservers)
	newNameserversLen := len(newNameservers)

	var emailType *string

	if emailTypeRaw, ok := data.GetOk("email_type"); ok {
		emailTypeString := emailTypeRaw.(string)
		emailType = &emailTypeString
	}

	nsResponse, err := client.DomainsDNS.GetList(domain)
	if err != nil {
		return diag.FromErr(err)
	}

	// If the previous state contains nameservers, but the new one does not contain,
	// then reset nameservers before applying records.
	// This case is possible when user removed nameservers lines and pasted records, so before applying records,
	// we must reset nameservers to defaults, otherwise we will face API exception
	if (oldNameserversLen != 0 && newNameserversLen == 0) ||
		// This condition resolves the issue if a user set up records on TF file, but in fact, manually enabled custom DNS.
		// Before applying records, we have to set default DNS
		(!*nsResponse.DomainDNSGetListResult.IsUsingOurDNS && newNameserversLen == 0) {
		_, err := client.DomainsDNS.SetDefault(domain)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if newRecordsLen != 0 || oldRecordsLen != 0 {
		diags := createRecordsOverwrite(ctx, domain, emailType, newRecords, client)
		if diags.HasError() {
			return diags
		}
	}

	if newNameserversLen != 0 {
		diags := createNameserversOverwrite(ctx, domain, convertInterfacesToString(newNameservers), client)
		if diags.HasError() {
			return diags
		}
	}

	// If user wants to control email type only while records & nameservers are absent,
	// then we have to update just an email status
	if emailType != nil && oldNameserversLen == 0 && newNameserversLen == 0 && oldRecordsLen == 0 && newRecordsLen == 0 {

		diags := createRecordsOverwrite(ctx, domain, emailType, []interface{}{}, client)
		if diags.HasError() {
			return diags
		}
	}

	// For overwrite mode, when no nameservers and records, and email type is not set, then we have to reset it to NONE
	if emailType == nil && oldNameserversLen == 0 && newNameserversLen == 0 && oldRecordsLen == 0 && newRecordsLen == 0 {
		diags := createRecordsOverwrite(ctx, domain, nil, []interface{}{}, client)
		if diags.HasError() {
			return diags
		}
	}

	return nil
}

func resourceRecordDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log(ctx, "resourceRecordDelete!!!!!!!!!!!!")

	client := meta.(*namecheap.Client)

	domain := strings.ToLower(data.Get("domain").(string))

	var records []interface{}
	var nameservers []interface{}

	if recordsRaw, ok := data.GetOk("record"); ok {
		records = recordsRaw.(*schema.Set).List()
	}

	if nameserversRaw, ok := data.GetOk("nameservers"); ok {
		nameservers = nameserversRaw.(*schema.Set).List()
	}

	recordsLen := len(records)
	nameserversLen := len(nameservers)

	if recordsLen != 0 {
		return deleteRecordsOverwrite(domain, client)
	}

	if nameserversLen != 0 {
		return deleteNameserversOverwrite(ctx, domain, client)
	}

	return nil
}
