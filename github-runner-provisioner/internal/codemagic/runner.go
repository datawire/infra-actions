package codemagic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type codemagicWorkflow struct {
	WorkflowID  string                       `json:"workflowId"`
	AppId       string                       `json:"appId"`
	Branch      string                       `json:"branch"`
	Environment map[string]map[string]string `json:"environment"`
}

func CreateMacM1Runner(ctx context.Context, owner string, repo string, token string, labels string, dryRun bool, codeMagicToken string) error {
	client := &http.Client{}

	data := codemagicWorkflow{
		WorkflowID: "github-runner",
		AppId:      "649493225428a76bc935a44b",
		Branch:     "main",
		Environment: map[string]map[string]string{
			"variables": {
				"GITHUB_REPO_OWNER":    owner,
				"GITHUB_REPO_NAME":     repo,
				"GITHUB_RUNNER_TOKEN":  token,
				"GITHUB_RUNNER_LABELS": labels,
			},
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(jsonData)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.codemagic.io/builds", buffer)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-auth-token", codeMagicToken)

	if dryRun {
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Error creating runner for ", owner, '/', repo, " Codemagic response: ", string(body))
		return errors.New("Failed to start CodeMagic build")
	}

	log.Info("Codemagic response: ", string(body))

	return nil
}
