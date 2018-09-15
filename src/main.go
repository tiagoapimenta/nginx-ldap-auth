package main

import (
	"flag"
	"fmt"
)

var config = flag.String("config", "/etc/nginx-ldap-auth/config.yaml", "Configuration file")

func main() {
	flag.Parse()
	fmt.Printf("Value of config: %s\n", *config)
}
