package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

func parseConfig() (string, *Config, error) {
	file := flag.String("config", "/etc/nginx-ldap-auth/config.yaml", "Configuration file")

	flag.Parse()

	data, err := ioutil.ReadFile(*file)
	if err != nil {
		return "", nil, fmt.Errorf("error on read file \"%s\": %v", *file, err)
	}

	config := Config{
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
			Success:     24 * time.Hour,
			Wrong:       5 * time.Minute,
			MaxAttempts: 0,
		},
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return "", nil, fmt.Errorf("Error on parse config: %v", err)
	}

	return *file, &config, nil
}
