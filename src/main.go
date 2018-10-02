package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

var (
	configFile = flag.String("config", "/etc/nginx-ldap-auth/config.yaml", "Configuration file")
	config     = Config{
		Web:     "0.0.0.0:5555",
		Path:    "/",
		Message: "LDAP Login",
		User: UserConfig{
			Filter: "(cn={0})",
		},
		Group: GroupConfig{
			Filter:    "(member={0})",
			GroupAttr: "cn",
		},
		Timeout: TimeoutConfig{
			Success: 24 * time.Hour,
			Wrong:   5 * time.Minute,
		},
	}
)

func main() {
	flag.Parse()

	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Error on read file \"%s\": %v\n", *configFile, err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error on parse config: %v\n", err)
	}

	fmt.Printf("Loaded config \"%s\".\n", *configFile)

	err = setupLDAP()
	if err != nil {
		log.Fatalf("Error on connect to LDAP: %v\n", err)
	}

	http.HandleFunc(config.Path, handler)

	fmt.Printf("Serving...\n")
	err = http.ListenAndServe(config.Web, nil)
	if err != nil {
		log.Fatalf("Error on start server: %v\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")

	if header != "" {
		auth := strings.SplitN(header, " ", 2)

		if len(auth) == 2 && auth[0] == "Basic" {
			decoded, err := base64.StdEncoding.DecodeString(auth[1])
			if err == nil {
				secret := strings.SplitN(string(decoded), ":", 2)

				if len(secret) == 2 && validate(secret[0], secret[1]) {
					// TODO: match by header, e.g: X-Original-URL X-Original-Method X-Sent-From X-Auth-Request-Redirect

					w.WriteHeader(http.StatusOK)
					return
				}
			} else {
				log.Printf("Error decode basic auth: %v\n", err)
			}
		}
	}

	w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", config.Message))
	w.WriteHeader(http.StatusUnauthorized)
}

func validate(username, password string) bool {
	ok, found := getCache(username, password)
	if found {
		return ok
	}

	ok, err := ldapLogin(username, password)
	if err != nil {
		log.Printf("Could not validade user %s: %v\n", username, err)
		return false
	}

	putCache(username, password, ok)
	return ok
}
