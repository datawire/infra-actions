package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ProvisioningError int

const (
	ErrorBadRequest ProvisioningError = iota
	ErrorInvalidAuthentication
	ErrorInvalidPayload
	ErrorUnknownAction
	ErrorUnknownRunnerLabel
	ErrorRunnerCreation
)

func (s ProvisioningError) String() string {
	switch s {
	case ErrorBadRequest:
		return "bad_request"
	case ErrorInvalidAuthentication:
		return "authentication_error"
	case ErrorInvalidPayload:
		return "invalid_request"
	case ErrorUnknownAction:
		return "unknown_workflow_action"
	case ErrorUnknownRunnerLabel:
		return "unknown_runner_label"
	case ErrorRunnerCreation:
		return "runner_creation_error"
	}
	return "unknown_error"
}

var ActionRunnerRuntime = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Subsystem: "action_runner",
	Name:      "runtime",
	Help:      "How long has an action runner been up."}, []string{"label", "instance_id"})

var RunnerProvisioningErrors = promauto.NewCounterVec(prometheus.CounterOpts{
	Subsystem: "action_runner",
	Name:      "provisioning_errors",
	Help:      "Errors managing runners on AWS."}, []string{"error", "runner_label"})
