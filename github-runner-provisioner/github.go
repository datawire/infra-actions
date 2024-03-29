package main

import (
	"context"

	"github.com/google/go-github/v48/github"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

func getGitHubAPIClient(ctx context.Context) *github.Client {
	cfg := NewConfig()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func getGitHubRunnerToken(ctx context.Context, owner string, repo string, dryRun bool) (token string, err error) {
	if dryRun {
		return "FAKE_TOKEN", nil
	}
	client := getGitHubAPIClient(ctx)
	registrationToken, _, err := client.Actions.CreateRegistrationToken(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	return *registrationToken.Token, nil
}

func getGitHubRunners(ctx context.Context, owner string, repo string) (*github.Runners, error) {
	client := getGitHubAPIClient(ctx)
	opts := &github.ListOptions{}
	runners, _, err := client.Actions.ListRunners(ctx, owner, repo, opts)
	if err != nil {
		return nil, err
	}
	return runners, nil
}

func isRunnerAvailable(ctx context.Context, owner string, repo string, labels []string) (bool, error) {
	runners, err := getGitHubRunners(ctx, owner, repo)
	if err != nil {
		return false, err
	}

	// check all runners registered to the repo
	for _, r := range runners.Runners {
		// if all labels were matched, check the availability
		if *r.Status != "online" || *r.Busy == true {
			continue
		}

		// check for label matches with each runner
		runnerLabelNames := []string{}
		for _, runnerLabel := range r.Labels {
			runnerLabelNames = append(runnerLabelNames, *runnerLabel.Name)
		}

		if labelsMatch(labels, runnerLabelNames) {
			return true, nil
		}
	}

	return false, nil
}

func labelsMatch(labels []string, runnerLabelNames []string) bool {
	for _, desiredLabel := range labels {
		if !slices.Contains(runnerLabelNames, desiredLabel) {
			return false
		}
	}
	return true
}
