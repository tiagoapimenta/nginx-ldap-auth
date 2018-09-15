package main

import "time"

type AuthConfig struct {
	BindDN string `yaml:"bindDN"`
	BindPW string `yaml:"bindPW"`
}

type UserConfig struct {
	BaseDN         string   `yaml:"baseDN"`
	Filter         string   `yaml:"filter"`
	UserAttr       string   `yaml:"userAttr"`
	RequiredGroups []string `yaml:"requiredGroups"`
}

type GroupConfig struct {
	BaseDN    string `yaml:"baseDN"`
	Filter    string `yaml:"filter"`
	UserAttr  string `yaml:"userAttr"`
	GroupAttr string `yaml:"groupAttr"`
}

type TimeoutConfig struct {
	Success time.Duration `yaml:"success"`
	Wrong   time.Duration `yaml:"wrong"`
}

type Config struct {
	Web     string        `yaml:"web"`
	Path    string        `yaml:"path"`
	Servers []string      `yaml:"servers"`
	Auth    AuthConfig    `yaml:"auth"`
	User    UserConfig    `yaml:"user"`
	Group   GroupConfig   `yaml:"group"`
	Timeout TimeoutConfig `yaml:"timeout"`
}
