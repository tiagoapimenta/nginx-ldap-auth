package main

import "time"

type AuthConfig struct {
	BindDN string `yaml:"bindDN"`
	BindPW string `yaml:"bindPW"`
}

type SearchConfig struct {
	BaseDN string `yaml:"baseDN"`
	Filter string `yaml:"filter"`
	Attr   string `yaml:"attr"`
}

type TimeoutConfig struct {
	Success time.Duration `yaml:"success"`
	Wrong   time.Duration `yaml:"wrong"`
}

type MatchConfig struct {
	Header string `yaml:"header"`
	Value  string `yaml:"value"`
	Regex  string `yaml:"regex"`
}

type PermissionConfig struct {
	Group string `yaml:"group"`
	User  string `yaml:"user"`
}

type RulesConfig struct {
	Match          []MatchConfig      `yaml:"match"`
	Allow          []PermissionConfig `yaml:"allow"`
	Deny           []PermissionConfig `yaml:"deny"`
	AllowAnonymous bool               `yaml:"allowAnonymous"`
}

type Config struct {
	Web     string        `yaml:"web"`
	Path    string        `yaml:"path"`
	Message string        `yaml:"message"`
	Servers []string      `yaml:"servers"`
	Auth    AuthConfig    `yaml:"auth"`
	User    SearchConfig  `yaml:"user"`
	Group   SearchConfig  `yaml:"group"`
	Timeout TimeoutConfig `yaml:"timeout"`
	Rules   []RulesConfig `yaml:"rules"`
}
