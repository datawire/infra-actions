package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/monitoring"
	"github.com/google/go-github/v48/github"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"strings"
)

func handleProvisioningRequest(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.String(), "/github-runner-provisioner") {
		message := fmt.Sprintf("URL %s is invalid", r.URL.String())
		http.Error(w, message, http.StatusBadRequest)
		log.Printf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorBadRequest.String(), "runner_label": ""}).Inc()
		return
	}

	if r.Method != http.MethodPost {
		message := "Only the POST method supported"
		http.Error(w, message, http.StatusBadRequest)
		log.Printf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorBadRequest.String(), "runner_label": ""}).Inc()
		return
	}

	payload, err := github.ValidatePayload(r, []byte(cfg.WebhookToken))
	if err != nil {
		message := "Webhook token invalid"
		http.Error(w, message, http.StatusUnauthorized)
		log.Printf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorInvalidAuthentication.String(), "runner_label": ""}).Inc()
		return
	}

	workflowJobEvent := github.WorkflowJobEvent{}
	err = json.Unmarshal(payload, &workflowJobEvent)
	if err != nil {
		message := fmt.Sprintf("Request is not a workflow job event: %v", err)
		http.Error(w, message, http.StatusBadRequest)
		log.Printf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorInvalidPayload.String(), "runner_label": ""}).Inc()
		return
	}

	if workflowJobEvent.Action == nil {
		message := "Workflow action is unknown"
		http.Error(w, message, http.StatusBadRequest)
		log.Printf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorUnknownAction.String(), "runner_label": ""}).Inc()
		return
	}

	if *workflowJobEvent.Action != "queued" {
		log.Printf("Ignoring GitHub event with action %s for repository %s", *workflowJobEvent.Action, *workflowJobEvent.Repo.Name)
		http.Error(w, "OK", http.StatusOK)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorUnknownAction.String(), "runner_label": ""}).Inc()
		return
	}

	var runnerFunction func(context.Context, string, string, bool) error
	var jobLabel string
	for _, label := range workflowJobEvent.WorkflowJob.Labels {
		if f, ok := runners[label]; ok {
			log.Printf("Job %s requested a runner with label %s\n", *workflowJobEvent.WorkflowJob.Name, label)
			runnerFunction = f
			jobLabel = label
			break
		}
	}

	if runnerFunction == nil {
		message := fmt.Sprintf("Workflow job didn't request a supported runner. Requested %v", workflowJobEvent.WorkflowJob.Labels)
		http.Error(w, message, http.StatusOK)
		log.Printf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorUnknownRunnerLabel.String(), "runner_label": ""}).Inc()
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
		message := fmt.Sprintf("Error creating %s runner for job %s [%s]: %v", jobLabel, *workflowJobEvent.WorkflowJob.Name, *workflowJobEvent.WorkflowJob.HTMLURL, err)
		http.Error(w, message, http.StatusInternalServerError)
		log.Printf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorRunnerCreation.String(), "runner_label": jobLabel}).Inc()
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
