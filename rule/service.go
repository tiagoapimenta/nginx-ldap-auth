package rule

import (
	"log"
	"sort"

	"github.com/tiagoapimenta/nginx-ldap-auth/data"
	"github.com/tiagoapimenta/nginx-ldap-auth/group"
	"github.com/tiagoapimenta/nginx-ldap-auth/user"
)

type Service struct {
	storage  *data.Storage
	user     *user.Service
	group    *group.Service
	required []string
}

func NewService(storage *data.Storage, userService *user.Service, groupService *group.Service, required []string) *Service {
	return &Service{
		storage:  storage,
		user:     userService,
		group:    groupService,
		required: required,
	}
}

func (p *Service) Validate(username, password string) bool {
	ok, found := p.storage.Get(username, password)
	if found {
		return ok
	}

	ok, err := p.validate(username, password)
	if err != nil {
		log.Printf("Could not validate user %s: %v\n", username, err)
		return false
	}

	p.storage.Put(username, password, ok)
	return ok
}

func (p *Service) validate(username, password string) (bool, error) {
	ok, id, err := p.user.Find(username)
	if !ok && err != nil {
		return false, err
	} else if err != nil {
		return false, nil
	}

	ok, err = p.user.Login(id, password)
	if !ok && err != nil {
		return false, err
	}

	if !ok || p.required == nil || len(p.required) == 0 {
		return err == nil, nil
	}

	groups, err := p.group.Find(id)
	if err != nil {
		return false, err
	}

	for _, group := range p.required {
		pos := sort.SearchStrings(groups, group)
		if pos >= len(groups) || groups[pos] != group {
			return false, nil
		}
	}

	return true, nil
}
