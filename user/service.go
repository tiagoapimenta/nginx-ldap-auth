package user

import (
	"strings"

	"github.com/tiagoapimenta/nginx-ldap-auth/ldap"

	gldap "github.com/go-ldap/ldap/v3"
)

type Service struct {
	pool   *ldap.Pool
	base   string
	filter string
}

func NewService(pool *ldap.Pool, base, filter string) *Service {
	return &Service{
		pool:   pool,
		base:   base,
		filter: filter,
	}
}

func (p *Service) Find(username string) (bool, string, error) {
	username = gldap.EscapeFilter(username)

	ok, id, _, err := p.pool.Search(
		p.base,
		strings.Replace(p.filter, "{0}", username, -1),
		"",
	)

	return ok, id, err
}

func (p *Service) Login(id, password string) (bool, error) {
	return p.pool.Validate(id, password)
}
