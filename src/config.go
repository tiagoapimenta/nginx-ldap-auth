package main

import "time"

type AuthConfig struct {
	BindDN string `yaml:"bindDN"`
	BindPW string `yaml:"bindPW"`
}

type UserConfig struct {
	BindDN         string   `yaml:"bindDN"`
	Filter         string   `yaml:"filter"`
	UserAttr       string   `yaml:"userAttr"`
	RequiredGroups []string `yaml:"requiredGroups"`
}

type GroupConfig struct {
	BindDN    string `yaml:"bindDN"`
	Filter    string `yaml:"filter"`
	UserAttr  string `yaml:"userAttr"`
	GroupAttr string `yaml:"member"`
}

type TimeoutConfig struct {
	Success time.Duration `yaml:"success"`
	Wrong   time.Duration `yaml:"wrong"`
}

type Config struct {
	Web     string        `yaml:"web"`
	Servers []string      `yaml:"servers"`
	Auth    AuthConfig    `yaml:"auth"`
	User    UserConfig    `yaml:"user"`
	Group   GroupConfig   `yaml:"group"`
	Timeout TimeoutConfig `yaml:"timeout"`
}
