package ldap

import (
	"fmt"
	"sort"

	ldap "github.com/go-ldap/ldap/v3"
)

func (p *Pool) Search(base, filter string, attr string) (bool, string, []string, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	err := p.auth()
	if err != nil {
		return false, "", nil, err
	}

	var list []string = nil
	if attr != "" {
		list = []string{attr}
	}

	var res *ldap.SearchResult
	_, err = p.networkJail(func() error {
		res, err = p.conn.Search(ldap.NewSearchRequest(
			base,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			filter,
			list,
			nil,
		))
		return err
	})

	if err != nil {
		return false, "", nil, err
	}

	if res == nil || len(res.Entries) == 0 {
		return true, "", nil, fmt.Errorf("No results for %s filter %s", base, filter)
	}

	if attr == "" && len(res.Entries) > 1 {
		return true, "", nil, fmt.Errorf("Too many results for %s filter %s", base, filter)
	}

	var result []string = nil
	if attr != "" {
		result = []string{}
		for _, entry := range res.Entries {
			result = append(result, entry.GetAttributeValue(attr))
		}
		sort.Strings(result)
	}

	return true, res.Entries[0].DN, result, nil
}
