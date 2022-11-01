package main

type Config struct {
	GithubToken  string `required:"true" envconfig:"GITHUB_TOKEN"`
	WebhookToken string `required:"true" envconfig:"WEBHOOK_TOKEN"`
}

var config *Config
