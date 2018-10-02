package main

import "time"

type AuthConfig struct {
	BindDN string `yaml:"bindDN"`
	BindPW string `yaml:"bindPW"`
}

type UserConfig struct {
	BaseDN         string   `yaml:"baseDN"`
	Filter         string   `yaml:"filter"`
	RequiredGroups []string `yaml:"requiredGroups"`
}

type GroupConfig struct {
	BaseDN    string `yaml:"baseDN"`
	GroupAttr string `yaml:"groupAttr"`
	Filter    string `yaml:"filter"`
}

type TimeoutConfig struct {
	Success time.Duration `yaml:"success"`
	Wrong   time.Duration `yaml:"wrong"`
}

type Config struct {
	Web     string        `yaml:"web"`
	Path    string        `yaml:"path"`
	Message string        `yaml:"message"`
	Servers []string      `yaml:"servers"`
	Auth    AuthConfig    `yaml:"auth"`
	User    UserConfig    `yaml:"user"`
	Group   GroupConfig   `yaml:"group"`
	Timeout TimeoutConfig `yaml:"timeout"`
}
