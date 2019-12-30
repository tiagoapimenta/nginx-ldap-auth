package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"

	ldap "github.com/go-ldap/ldap/v3"
)

func (p *Pool) Connect() error {
	if p.url == "" {
		return errors.New("No LDAP server available!")
	}

	if p.port == 0 {
		return fmt.Errorf("Unable to determine schema or port for \"%s\"", p.url)
	}

	if p.conn != nil {
		p.conn.Close()
	}

	address := fmt.Sprintf("%s:%d", p.url, p.port)
	if p.ssl {
		conn, err := ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}
		p.conn = conn
	} else {
		conn, err := ldap.Dial("tcp", address)
		if err != nil {
			return err
		}
		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			log.Printf("It was not possble to start TLS, falling back to plain: %v.\n", err)
			conn.Close()
			conn, err = ldap.Dial("tcp", address)
			if err != nil {
				return err
			}
		}
		p.conn = conn
	}

	p.admin = false

	return p.auth()
}
