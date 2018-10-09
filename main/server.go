package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/tiagoapimenta/nginx-ldap-auth/rule"
)

func startServer(service *rule.Service, server, path, message string) error {
	realm := fmt.Sprintf("Basic realm=\"%s\"", message)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header != "" {
			auth := strings.SplitN(header, " ", 2)

			if len(auth) == 2 && auth[0] == "Basic" {
				decoded, err := base64.StdEncoding.DecodeString(auth[1])
				if err == nil {
					secret := strings.SplitN(string(decoded), ":", 2)

					if len(secret) == 2 && service.Validate(secret[0], secret[1]) {
						w.WriteHeader(http.StatusOK)
						return
					}
				} else {
					log.Printf("Error decode basic auth: %v\n", err)
				}
			}
		}

		w.Header().Set("WWW-Authenticate", realm)
		w.WriteHeader(http.StatusUnauthorized)
	})

	return http.ListenAndServe(server, nil)
}
