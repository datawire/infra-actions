package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v48/github"
	"golang.org/x/exp/slices"
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

	//workflowJobEvent := github.WorkflowJobEvent{}
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	if len(r.PostForm) == 0 {
		http.Error(w, "Request contains no data", http.StatusBadRequest)
		return
	}

	formData, ok := r.PostForm["payload"]
	if !ok {
		http.Error(w, "Request has no payload", http.StatusBadRequest)
	}

	workflowJobEvent := github.WorkflowJobEvent{}
	err := json.Unmarshal([]byte(formData[0]), &workflowJobEvent)
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

	if !slices.Contains(workflowJobEvent.WorkflowJob.Labels, "macOS-arm64") {
		http.Error(w, fmt.Sprintf("Only runners of type macOS-arm64 are supported. Got %v", workflowJobEvent.WorkflowJob.Labels), http.StatusOK)
		return
	}

	log.Printf("Job %s requested a Mac M1 runner\n", *workflowJobEvent.Repo.Name)

	dryRun := len(r.Form["dry-run"]) > 0 && r.Form["dry-run"][0] == "true"
	if err := createMacM1Runner(r.Context(), *workflowJobEvent.Repo.Owner.Login, *workflowJobEvent.Repo.Name, dryRun); err != nil {
		log.Printf("Error creating Mac M1 runner for job %s [%s]: %v", *workflowJobEvent.WorkflowJob.Name, *workflowJobEvent.WorkflowJob.HTMLURL, err)
		http.Error(w, fmt.Sprintf("Error creating Mac M1 runner: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Mac M1 runner has been scheduled for job %s\n", *workflowJobEvent.Repo.Name)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("Error sending HTTP response: %v", err)
	}
}

func handleHealthCheckRequest(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "OK", http.StatusOK)
}
