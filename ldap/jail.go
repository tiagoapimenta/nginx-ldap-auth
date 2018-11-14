package ldap

import (
	"log"

	ldap "gopkg.in/ldap.v2"
)

func (p *Pool) networkJail(f func() error) error {
	err := f()
	if err != nil && ldap.IsErrorWithCode(err, ldap.ErrorNetwork) {
		log.Printf("Network problem, trying to reconnect once: %v.\n", err)
		err = p.Connect()
		if err != nil {
			return err
		}
		err = f()
	}
	return err
}
