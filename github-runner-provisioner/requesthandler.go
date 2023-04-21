package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/datawire/infra-actions/github-runner-provisioner/internal/monitoring"
	"github.com/google/go-github/v48/github"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func setupLogFields(r *http.Request, status int, requestTime time.Time) log.Fields {
	return log.Fields{
		"httpRequest": log.Fields{
			"requestMethod": r.Method,
			"requestUrl":    r.URL.String(),
			"requestSize":   r.ContentLength,
			"status":        status,
			"userAgent":     r.UserAgent(),
			"remoteIp":      r.RemoteAddr,
			"referer":       r.Referer(),
			"latency":       time.Since(requestTime).String(),
		},
	}
}

func handleProvisioningRequest(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now()
	log.WithFields(setupLogFields(r, 200, requestTime)).Info("Request received")

	if !strings.HasPrefix(r.URL.String(), "/github-runner-provisioner") {
		message := fmt.Sprintf("URL %s is invalid", r.URL.String())
		http.Error(w, message, http.StatusBadRequest)
		log.WithFields(setupLogFields(r, http.StatusBadRequest, requestTime)).Errorf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorBadRequest.String(), "runner_label": "", "repo": ""}).Inc()
		return
	}

	if r.Method != http.MethodPost {
		message := "Only the POST method supported"
		http.Error(w, message, http.StatusMethodNotAllowed)
		log.WithFields(setupLogFields(r, http.StatusMethodNotAllowed, requestTime)).Errorf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorBadRequest.String(), "runner_label": "", "repo": ""}).Inc()
		return
	}

	payload, err := github.ValidatePayload(r, []byte(cfg.WebhookToken))
	if err != nil {
		message := "Webhook token invalid"
		http.Error(w, message, http.StatusUnauthorized)
		log.WithFields(setupLogFields(r, http.StatusUnauthorized, requestTime)).Errorf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorInvalidAuthentication.String(), "runner_label": "", "repo": ""}).Inc()
		return
	}

	workflowJobEvent := github.WorkflowJobEvent{}
	err = json.Unmarshal(payload, &workflowJobEvent)
	if err != nil {
		message := fmt.Sprintf("Request is not a workflow job event: %v", err)
		http.Error(w, message, http.StatusBadRequest)
		log.WithFields(setupLogFields(r, http.StatusBadRequest, requestTime)).Warningf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"warning": monitoring.ErrorInvalidPayload.String(), "runner_label": "", "repo": ""}).Inc()
		return
	}

	if workflowJobEvent.Action == nil {
		message := "Workflow action is unknown"
		http.Error(w, message, http.StatusBadRequest)
		log.WithFields(setupLogFields(r, http.StatusBadRequest, requestTime)).Errorf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorUnknownAction.String(),
			"runner_label": "", "repo": *workflowJobEvent.Repo.Name}).Inc()
		return
	}

	if *workflowJobEvent.Action != "queued" {
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
		log.WithFields(setupLogFields(r, http.StatusOK, requestTime)).Infof("Ignoring GitHub event with action %s for repository %s", *workflowJobEvent.Action, *workflowJobEvent.Repo.Name)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"info": monitoring.ErrorUnknownAction.String(),
			"runner_label": "", "repo": *workflowJobEvent.Repo.Name}).Inc()
		return
	}

	var runnerFunction func(context.Context, string, string, bool) error
	var jobLabel string
	for _, label := range workflowJobEvent.WorkflowJob.Labels {
		if f, ok := runners[label]; ok {
			log.Infof("Job %s requested a runner with label %s", *workflowJobEvent.WorkflowJob.Name, label)
			runnerFunction = f
			jobLabel = label
			break
		}
	}

	if runnerFunction == nil {
		message := fmt.Sprintf("Workflow job didn't request a supported runner. Requested %v", workflowJobEvent.WorkflowJob.Labels)
		http.Error(w, message, http.StatusOK)
		log.WithFields(setupLogFields(r, http.StatusOK, requestTime)).Infof(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"info": monitoring.ErrorUnknownRunnerLabel.String(),
			"runner_label": jobLabel, "repo": *workflowJobEvent.Repo.Name}).Inc()
		return
	}

	log.Infof("Job in %s repo requested a %s runner", *workflowJobEvent.Repo.Name, jobLabel)

	runnerLabels := []string{0: jobLabel}
	isAvailable, err := isRunnerAvailable(r.Context(), *workflowJobEvent.Repo.Owner.Login, *workflowJobEvent.Repo.Name, runnerLabels)
	if err != nil {
		message := fmt.Sprintf("Error checking if runner is available: %v", err)
		http.Error(w, message, http.StatusInternalServerError)
		log.WithFields(setupLogFields(r, http.StatusInternalServerError, requestTime)).Errorf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorCheckingAvailableRunners.String(),
			"runner_label": jobLabel, "repo": *workflowJobEvent.Repo.Name}).Inc()
		return
	}

	if isAvailable {
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
		log.WithFields(setupLogFields(r, http.StatusOK, requestTime)).Infof("%s runner already available. No scaling action required.", jobLabel)
		return
	}

	dryRun := len(r.Form["dry-run"]) > 0 && r.Form["dry-run"][0] == "true"
	if err := runnerFunction(r.Context(), *workflowJobEvent.Repo.Owner.Login, *workflowJobEvent.Repo.Name, dryRun); err != nil {
		message := fmt.Sprintf("Error creating %s runner for job %s [%s]: %v", jobLabel, *workflowJobEvent.WorkflowJob.Name, *workflowJobEvent.WorkflowJob.HTMLURL, err)
		http.Error(w, message, http.StatusInternalServerError)
		log.WithFields(setupLogFields(r, http.StatusInternalServerError, requestTime)).Errorf(message)

		monitoring.RunnerProvisioningErrors.With(prometheus.Labels{"error": monitoring.ErrorRunnerCreation.String(),
			"runner_label": jobLabel, "repo": *workflowJobEvent.Repo.Name}).Inc()
		return
	}

	http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
	log.WithFields(setupLogFields(r, http.StatusOK, requestTime)).Infof("%s runner has been scheduled for job %s", jobLabel, *workflowJobEvent.Repo.Name)
}

func handleHealthCheckRequest(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "OK", http.StatusOK)
}
