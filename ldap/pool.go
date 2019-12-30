package ldap

import (
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
)

type Pool struct {
	url      string
	port     int
	ssl      bool
	username string
	password string
	conn     *ldap.Conn
	admin    bool
	lock     sync.Mutex
}

func NewPool(servers []string, username, password string) *Pool {
	url := ""
	port := 0
	schema := "auto"

	size := len(servers)
	if size != 0 {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		server := servers[r.Intn(size)]

		url = server
		if strings.HasPrefix(url, "ldaps:") {
			url = strings.TrimPrefix(strings.TrimPrefix(url, "ldaps:"), "//")
			schema = "ldaps"
			port = 636
		} else if strings.HasPrefix(url, "ldap:") {
			url = strings.TrimPrefix(strings.TrimPrefix(url, "ldap:"), "//")
			schema = "ldap"
			port = 389
		}

		portExp := regexp.MustCompile(`:[0-9]+$`)
		if portExp.MatchString(url) {
			str := portExp.FindString(url)

			number, err := strconv.Atoi(str[1:])
			if err == nil {
				port = number
				url = strings.TrimSuffix(url, str)
			} else {
				log.Printf("Error on parse port of \"%s\": %v\n", server, err)
			}
		}

		if schema == "auto" {
			if port == 636 {
				schema = "ldaps"
			} else if port == 389 {
				schema = "ldap"
			} else {
				port = 0
			}
		}
	}

	return &Pool{
		url:      url,
		port:     port,
		ssl:      schema == "ldaps",
		username: username,
		password: password,
	}
}
