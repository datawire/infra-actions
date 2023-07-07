package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	GithubToken    string `required:"true" envconfig:"GITHUB_TOKEN"`
	WebhookToken   string `required:"true" envconfig:"WEBHOOK_TOKEN"`
	CodeMagicToken string `required:"true" envconfig:"CODEMAGIC_TOKEN"`
	UseCodeMagic   bool   `default:"true"  envconfig:"USE_CODEMAGIC"`
}

func NewConfig() *Config {
	cfg := &Config{}

	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("Error loading environment configuration: %v", err)
	}

	return cfg
}
