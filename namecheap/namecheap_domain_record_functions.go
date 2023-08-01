package namecheap_provider

import (
	"context"
	"fmt"
	"github.com/agent-tao/go-namecheap-sdk/v2/namecheap"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strings"
)

// createNameserversOverwrite force overwrites the nameservers
func createNameserversOverwrite(ctx context.Context, domain string, nameservers []string, client *namecheap.Client) diag.Diagnostics {
	log(ctx, "createNameserversOverwrite!!!!!!!")

	_, err := client.DomainsDNS.SetCustom(domain, nameservers)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// readNameserversOverwrite returns remote real nameservers
func readNameserversOverwrite(ctx context.Context, domain string, client *namecheap.Client) (*[]string, diag.Diagnostics) {
	log(ctx, "readNameserversOverwrite!!!!!!!")

	nsResponse, err := client.DomainsDNS.GetList(domain)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if *nsResponse.DomainDNSGetListResult.IsUsingOurDNS || nsResponse.DomainDNSGetListResult.Nameservers == nil {
		return &[]string{}, nil
	} else {
		return nsResponse.DomainDNSGetListResult.Nameservers, nil
	}
}

// deleteNameserversOverwrite resets nameservers settings to default (set default Namecheap's nameservers)
func deleteNameserversOverwrite(ctx context.Context, domain string, client *namecheap.Client) diag.Diagnostics {
	log(ctx, "deleteNameserversOverwrite!!!!!!!")

	_, err := client.DomainsDNS.SetDefault(domain)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// createRecordsOverwrite overwrites existing records with provided new ones
func createRecordsOverwrite(ctx context.Context, domain string, emailType *string, records []interface{}, client *namecheap.Client) diag.Diagnostics {
	log(ctx, "createRecordsOverwrite!!!!!!!!!!!!")
	domainRecords := convertRecordTypeSetToDomainRecords(&records)

	emailTypeValue := namecheap.String(namecheap.EmailTypeNone)
	if emailType != nil {
		emailTypeValue = emailType
	}

	_, err := client.DomainsDNS.SetHosts(&namecheap.DomainsDNSSetHostsArgs{
		Domain:    &domain,
		Records:   domainRecords,
		EmailType: emailTypeValue,
		Flag:      nil,
		Tag:       nil,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

// readRecordsOverwrite returns the records that are exist on Namecheap
// NOTE: method has address fix. Refer to getFixedAddressOfRecord
func readRecordsOverwrite(domain string, currentRecords []interface{}, client *namecheap.Client) (*[]map[string]interface{}, *string, diag.Diagnostics) {
	remoteRecordsResponse, err := client.DomainsDNS.GetHosts(domain)
	if err != nil {
		return nil, nil, diag.FromErr(err)
	}

	currentRecordsConverted := convertRecordTypeSetToDomainRecords(&currentRecords)

	var remoteRecords []map[string]interface{}

	if remoteRecordsResponse.DomainDNSGetHostsResult.Hosts != nil {
		for _, remoteRecord := range *remoteRecordsResponse.DomainDNSGetHostsResult.Hosts {
			remoteRecordHash := hashRecord(*remoteRecord.Name, *remoteRecord.Type, *remoteRecord.Address)

			for _, currentRecord := range *currentRecordsConverted {
				currentRecordAddressFixed, err := getFixedAddressOfRecord(&currentRecord)
				if err != nil {
					return nil, nil, diag.FromErr(err)
				}

				currentRecordHash := hashRecord(*currentRecord.HostName, *currentRecord.RecordType, *currentRecordAddressFixed)

				if currentRecordHash == remoteRecordHash {
					*remoteRecord.Address = *currentRecord.Address
					break
				}

			}

			remoteRecords = append(remoteRecords, *convertDomainRecordDetailedToTypeSetRecord(&remoteRecord))
		}
	}

	return &remoteRecords, remoteRecordsResponse.DomainDNSGetHostsResult.EmailType, nil
}

// deleteRecordsOverwrite removes all records
func deleteRecordsOverwrite(domain string, client *namecheap.Client) diag.Diagnostics {
	var records []namecheap.DomainsDNSHostRecord

	_, err := client.DomainsDNS.SetHosts(&namecheap.DomainsDNSSetHostsArgs{
		Domain:    &domain,
		Records:   &records,
		EmailType: namecheap.String(namecheap.EmailTypeNone),
		Flag:      nil,
		Tag:       nil,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// hashRecord creates a hash for record by hostname, recordType, and address
func hashRecord(hostname string, recordType string, address string) string {
	return fmt.Sprintf("[%s:%s:%s]", hostname, recordType, address)
}

func convertRecordTypeSetToDomainRecords(records *[]interface{}) *[]namecheap.DomainsDNSHostRecord {
	var mappedRecords []namecheap.DomainsDNSHostRecord

	for _, _record := range *records {
		record := _record.(map[string]interface{})

		hostNameVal := record["hostname"].(string)
		typeVal := record["type"].(string)
		addressVal := record["address"].(string)
		mxPrefVal := record["mx_pref"].(int)
		ttlVal := record["ttl"].(int)

		domainRecord := namecheap.DomainsDNSHostRecord{
			HostName:   namecheap.String(hostNameVal),
			RecordType: namecheap.String(typeVal),
			Address:    namecheap.String(addressVal),
			MXPref:     namecheap.UInt8(uint8(mxPrefVal)),
			TTL:        namecheap.Int(ttlVal),
		}

		mappedRecords = append(mappedRecords, domainRecord)
	}

	return &mappedRecords
}

func convertDomainRecordDetailedToTypeSetRecord(record *namecheap.DomainsDNSHostRecordDetailed) *map[string]interface{} {
	return &map[string]interface{}{
		"hostname": *record.Name,
		"type":     *record.Type,
		"address":  *record.Address,
		"mx_pref":  *record.MXPref,
		"ttl":      *record.TTL,
	}
}

func convertInterfacesToString(stringsRaw []interface{}) []string {
	var stringList []string
	for _, stringRaw := range stringsRaw {
		stringList = append(stringList, stringRaw.(string))
	}
	return stringList
}

func fixCAAIodefAddressValue(address *string) (*string, error) {
	addressValues := strings.Split(strings.TrimSpace(*address), " ")
	var addressValuesFixed []string

	for _, value := range addressValues {
		fixedValue := strings.TrimSpace(value)
		if len(fixedValue) != 0 {
			addressValuesFixed = append(addressValuesFixed, fixedValue)
		}
	}

	if len(addressValuesFixed) != 3 {
		return nil, fmt.Errorf(`Invalid value "%s"`, *address)
	}

	hasPrefixQuote := strings.HasPrefix(addressValuesFixed[2], `"`)
	hasSuffixQuite := strings.HasSuffix(addressValuesFixed[2], `"`)

	if !hasPrefixQuote && !hasSuffixQuite {
		addressValuesFixed[2] = fmt.Sprintf(`"%s"`, addressValuesFixed[2])
	} else if !hasPrefixQuote || !hasSuffixQuite {
		return nil, fmt.Errorf(`Invalid value "%s"`, *address)
	}

	addressNew := strings.Join(addressValuesFixed, " ")
	return &addressNew, nil
}

func fixAddressEndWithDot(address *string) *string {
	if !strings.HasSuffix(*address, ".") {
		return namecheap.String(*address + ".")
	}
	return address
}

// getFixedAddressOfRecord check the record type and return the fixed address with either dot suffix or quotes around domain name
// The following addresses should be returned:
// - for CNAME, ALIAS, NS, MX records, if the address has been provided without dot suffix, then it will be added
// - for CAA records with iodef key word, if no quotes wrapping the domain, then the quotes will be added
// - for other cases the method will just return the address equal to input one
func getFixedAddressOfRecord(record *namecheap.DomainsDNSHostRecord) (*string, error) {
	if *record.RecordType == namecheap.RecordTypeCNAME ||
		*record.RecordType == namecheap.RecordTypeAlias ||
		*record.RecordType == namecheap.RecordTypeNS ||
		*record.RecordType == namecheap.RecordTypeMX {
		return fixAddressEndWithDot(record.Address), nil
	}

	if *record.RecordType == namecheap.RecordTypeCAA && strings.Contains(*record.Address, "iodef") {
		return fixCAAIodefAddressValue(record.Address)
	}

	return record.Address, nil
}

// filterDefaultParkingRecords filters out default parking records
func filterDefaultParkingRecords(records *[]namecheap.DomainsDNSHostRecordDetailed, domain *string) *[]namecheap.DomainsDNSHostRecordDetailed {
	var filteredRecords []namecheap.DomainsDNSHostRecordDetailed

	for _, record := range *records {
		if (*record.Type == namecheap.RecordTypeCNAME && *record.Name == "www" && *record.Address == "parkingpage.namecheap.com.") ||
			(*record.Type == namecheap.RecordTypeURL && *record.Name == "@" && strings.HasPrefix(*record.Address, "http://www."+*domain)) {
			continue
		}
		filteredRecords = append(filteredRecords, record)
	}

	return &filteredRecords
}

// stringifyNCRecord returns a string with hostname, record type and address of the record
// This function mostly serves to print error details for user
func stringifyNCRecord(record *namecheap.DomainsDNSHostRecord) string {
	return fmt.Sprintf("{hostname = %s, type = %s, address = %s}", *record.HostName, *record.RecordType, *record.Address)
}

// resolveEmailType resolves an emailType for the case when no emailType provided by terraform configuration,
// but we have an old emailType value extracted from read response
// The main purpose is to prevent set up MX/MXE email type when after manipulation no MX/MXE records available
// This function resolves a bug when we have removed MX/MXE record without reset of emailType, then trying to remove non-MX* record
func resolveEmailType(records *[]namecheap.DomainsDNSHostRecord, emailType *string) *string {
	if *emailType != namecheap.EmailTypeMXE && *emailType != namecheap.EmailTypeMX {
		return emailType
	}

	foundMX := false
	foundMXE := false

	for _, record := range *records {
		if *record.RecordType == namecheap.RecordTypeMX {
			foundMX = true
		} else if *record.RecordType == namecheap.RecordTypeMXE {
			foundMXE = true
		}
	}

	if *emailType == namecheap.EmailTypeMX && !foundMX ||
		*emailType == namecheap.EmailTypeMXE && !foundMXE {
		return namecheap.String(namecheap.EmailTypeNone)
	}

	return emailType
}

func createDomainIfNonexist(ctx context.Context, domain string, client *namecheap.Client) diag.Diagnostics {
	//get domain info
	_, err := client.Domains.GetInfo(domain)

	//if domain does not exist, then create
	if err != nil {
		log(ctx, "Can not Get Domain Info:%s", domain)

		//log.Println("Can not Get Domain Info, Creating:%s", domain)
		resp, err := client.Domains.DomainsAvailable(domain)
		if err == nil && *resp.Result.Available == true {
			// no err and available, create
			log(ctx, "Can not Get Domain Info, Creating %s", domain)
			_, err = client.Domains.DomainsCreate(domain, _info)

			if err != nil {
				log(ctx, "create domain %s failed, exit", domain)
				log(ctx, "reason:", err.Error())
				return diag.Errorf("create domain failed", domain)
			}

		} else {
			log(ctx, "domain %s is not available, exiting!", domain)
			return diag.Errorf("domain is not available to register, you need to change to another domain", domain)
		}
	} else {
		//skip, do nothing
		tflog.Info(ctx, "Domain %s exist, then do record config", domain)
	}
	return nil

}

func log(ctx context.Context, format string, a ...any) {
	msg := fmt.Sprintf(format, a)
	tflog.Info(ctx, msg)

}
