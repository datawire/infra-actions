package main

import (
	"context"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

func getGitHubRunnerToken(ctx context.Context, owner string, repo string) (token string, err error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	registrationToken, _, err := client.Actions.CreateRegistrationToken(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	return *registrationToken.Token, nil
}
