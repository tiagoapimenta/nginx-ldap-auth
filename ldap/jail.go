package ldap

import (
	"log"

	ldap "github.com/go-ldap/ldap/v3"
)

func (p *Pool) networkJail(f func() error) (bool, error) {
	err := f()
	if err != nil && ldap.IsErrorWithCode(err, ldap.ErrorNetwork) {
		log.Printf("Network problem, trying to reconnect once: %v.\n", err)
		err = p.Connect()
		if err != nil {
			return false, err
		}
		err = f()
		if err != nil && ldap.IsErrorWithCode(err, ldap.ErrorNetwork) {
			return false, err
		}
	}
	return true, err
}
