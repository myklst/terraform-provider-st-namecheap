package namecheap_provider

import (
	"github.com/agent-tao/go-namecheap-sdk/v2/namecheap"
)

// equalDomainRecord compares only Name, Type, Address, TTL, MXPref fields only
func equalDomainRecord(sRec *namecheap.DomainsDNSHostRecordDetailed, dRec *namecheap.DomainsDNSHostRecordDetailed) bool {
	return *sRec.Name == *dRec.Name &&
		*sRec.Type == *dRec.Type &&
		*sRec.Address == *dRec.Address &&
		*sRec.TTL == *dRec.TTL &&
		*sRec.MXPref == *dRec.MXPref
}
