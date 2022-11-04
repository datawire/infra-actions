package main

import (
	"context"
	"github.com/google/go-github/v48/github"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

func getGitHubAPIClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func getGitHubRunnerToken(ctx context.Context, owner string, repo string) (token string, err error) {
	client := getGitHubAPIClient(ctx)
	registrationToken, _, err := client.Actions.CreateRegistrationToken(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	return *registrationToken.Token, nil
}

func getGitHubRunners(ctx context.Context, owner string, repo string) *github.Runners {
	client := getGitHubAPIClient(ctx)
	opts := &github.ListOptions{}
	runners, _, err := client.Actions.ListRunners(ctx, owner, repo, opts)
	if err != nil {
		return nil
	}
	return runners
}

func isRunnerAvailable(ctx context.Context, owner string, repo string, labels []string) bool {
	runners := getGitHubRunners(ctx, owner, repo)

	for i := range runners.Runners {
		r := runners.Runners[i]
		var matches = []bool{}
		for j := range labels {
			matches = append(matches, false)
			for k := range r.Labels {
				if labels[j] == *r.Labels[k].Name {
					matches[j] = true
					break
				}
			}
		}
		if !slices.Contains(matches, false) {
			if *r.Status == "online" && *r.Busy == false {
				return true
			}
		}
	}

	return false
}
