package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v48/github"
	"log"
	"net/http"
	"strings"
)

func handleProvisioningRequest(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.String(), "/github-runner-provisioner") {
		http.Error(w, fmt.Sprintf("URL %s is invalid", r.URL.String()), http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only the POST method supported", http.StatusBadRequest)
		return
	}

	payload, err := github.ValidatePayload(r, []byte(cfg.WebhookToken))
	if err != nil {
		http.Error(w, "Webhook token invalid", http.StatusUnauthorized)
		return
	}

	workflowJobEvent := github.WorkflowJobEvent{}
	err = json.Unmarshal(payload, &workflowJobEvent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Request is not a workflow job event: %v", err), http.StatusBadRequest)
		return
	}

	if workflowJobEvent.Action == nil {
		http.Error(w, "Workflow action is unknown", http.StatusBadRequest)
		return
	}

	if *workflowJobEvent.Action != "queued" {
		log.Printf("Ignoring GitHub event with action %s.", *workflowJobEvent.Action)
		http.Error(w, "OK", http.StatusOK)
		return
	}

	var runnerFunction func(context.Context, string, string, bool) error
	var jobLabel string
	for _, label := range workflowJobEvent.WorkflowJob.Labels {
		if f, ok := runners[label]; ok {
			log.Printf("Job %s requested a runner with label %s\n", workflowJobEvent.WorkflowJob.Name, label)
			runnerFunction = f
			jobLabel = label
			break
		}
	}

	if runnerFunction == nil {
		http.Error(w, fmt.Sprintf("Workflow job didn't request a supported runner. Requested %v", workflowJobEvent.WorkflowJob.Labels), http.StatusOK)
		return
	}

	log.Printf("Job in %s repo requested a %s runner\n", *workflowJobEvent.Repo.Name, jobLabel)

	runnerLabels := []string{0: jobLabel}
	if isRunnerAvailable(r.Context(), *workflowJobEvent.Repo.Owner.Login, *workflowJobEvent.Repo.Name, runnerLabels) {
		log.Printf("%s runner already available. No action scaling action required.", jobLabel)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Error sending HTTP response: %v", err)
		}
	}

	dryRun := len(r.Form["dry-run"]) > 0 && r.Form["dry-run"][0] == "true"
	if err := runnerFunction(r.Context(), *workflowJobEvent.Repo.Owner.Login, *workflowJobEvent.Repo.Name, dryRun); err != nil {
		log.Printf("Error creating %s runner for job %s [%s]: %v", jobLabel, *workflowJobEvent.WorkflowJob.Name, *workflowJobEvent.WorkflowJob.HTMLURL, err)
		http.Error(w, fmt.Sprintf("Error creating %s runner: %v", jobLabel, err), http.StatusInternalServerError)
		return
	}

	log.Printf("%s runner has been scheduled for job %s\n", jobLabel, *workflowJobEvent.Repo.Name)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("Error sending HTTP response: %v", err)
	}
}

func handleHealthCheckRequest(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "OK", http.StatusOK)
}
