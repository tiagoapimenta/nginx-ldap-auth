package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	ldap "gopkg.in/ldap.v2"
)

var conn *ldap.Conn

func setupLDAP() error {
	size := len(config.Servers)
	if size == 0 {
		return fmt.Errorf("No LDAP server available on %+v", config)
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	portExp := regexp.MustCompile(`:[0-9]+$`)

	server := config.Servers[r.Intn(size)]
	url := server
	port := 0
	schema := "auto"

	if strings.HasPrefix(url, "ldaps:") {
		url = strings.TrimPrefix(strings.TrimPrefix(url, "ldaps:"), "//")
		schema = "ldaps"
		port = 636
	} else if strings.HasPrefix(url, "ldap:") {
		url = strings.TrimPrefix(strings.TrimPrefix(url, "ldap:"), "//")
		schema = "ldap"
		port = 389
	}

	var err error
	if portExp.MatchString(url) {
		str := portExp.FindString(url)
		port, err = strconv.Atoi(str[1:])
		if err != nil {
			return fmt.Errorf("Unable to parse url \"%s\", %v", server, err)
		}
		url = strings.TrimSuffix(url, str)
	}

	if schema == "auto" {
		if port == 636 {
			schema = "ldaps"
		} else if port == 389 {
			schema = "ldap"
		}
	}

	if schema == "auto" || port == 0 {
		return fmt.Errorf("Unable to determine schema or port for \"%s\"", server)
	}

	address := fmt.Sprintf("%s:%d", url, port)
	fmt.Printf("Connecting LDAP %s...\n", address)

	if schema == "ldaps" {
		conn, err = ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}
	} else {
		conn, err = ldap.Dial("tcp", address)
		if err != nil {
			return err
		}
		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}
	}

	if config.Auth.BindDN != "" || config.Auth.BindPW != "" {
		err = conn.Bind(config.Auth.BindDN, config.Auth.BindPW)
		if err != nil {
			return err
		}
	}
	return nil
}

func ldapLogin(username, password string) (bool, error) {
	// TODO: lock
	if config.Auth.BindDN != "" || config.Auth.BindPW != "" {
		err := conn.Bind(config.Auth.BindDN, config.Auth.BindPW)
		if err != nil {
			return false, err
		}
	}

	req := ldap.NewSearchRequest(
		config.User.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		strings.Replace(config.User.Filter, "{0}", username, -1),
		[]string{config.User.UserAttr},
		nil,
	)

	res, err := conn.Search(req)

	if err != nil {
		return false, err
	}

	if len(res.Entries) != 1 {
		return false, nil
	}

	err = conn.Bind(res.Entries[0].DN, password)
	return err == nil, nil
}
