package main

type Config struct {
	GithubToken string `required:"true" envconfig:"GITHUB_TOKEN"`
}

var config *Config
