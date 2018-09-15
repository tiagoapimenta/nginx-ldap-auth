package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	yaml "gopkg.in/yaml.v2"
)

var (
	configFile = flag.String("config", "/etc/nginx-ldap-auth/config.yaml", "Configuration file")
	config     = Config{
		Web:  "0.0.0.0:5555",
		Path: "/",
		User: UserConfig{
			UserAttr: "uid",
		},
		Group: GroupConfig{
			UserAttr:  "uid",
			GroupAttr: "member",
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

	fmt.Printf("Loaded config: %+v\n", config)

	err = setupLDAP()
	if err != nil {
		log.Fatalf("Error on connect to LDAP: %v\n", err)
	}

	http.HandleFunc(config.Path, handler)

	err = http.ListenAndServe(config.Web, nil)
	if err != nil {
		log.Fatalf("Error on start server: %v\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusOK)
}
