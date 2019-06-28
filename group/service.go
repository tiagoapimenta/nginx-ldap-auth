package group

import (
	"strings"

	"github.com/tiagoapimenta/nginx-ldap-auth/ldap"

	gldap "gopkg.in/ldap.v2"
)

type Service struct {
	pool   *ldap.Pool
	base   string
	filter string
	attr   string
}

func NewService(pool *ldap.Pool, base, filter, attr string) *Service {
	return &Service{
		pool:   pool,
		base:   base,
		filter: filter,
		attr:   attr,
	}
}

func (p *Service) Find(id string) ([]string, error) {
	id = gldap.EscapeFilter(id)

	ok, _, groups, err := p.pool.Search(
		p.base,
		strings.Replace(p.filter, "{0}", id, -1),
		p.attr,
	)

	if !ok && err != nil {
		return nil, err
	} else if err != nil {
		return []string{}, nil
	}

	return groups, nil
}
