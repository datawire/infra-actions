package main

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	GithubToken  string `required:"true" envconfig:"GITHUB_TOKEN"`
	WebhookToken string `required:"true" envconfig:"WEBHOOK_TOKEN"`
}

func NewConfig() *Config {
	cfg := &Config{}

	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("Error loading environment configuration: %v", err)
	}

	return cfg
}
