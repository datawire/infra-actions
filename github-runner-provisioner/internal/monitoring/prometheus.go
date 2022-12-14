package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var ActionRunnerRuntime = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Subsystem: "action_runner",
	Name:      "runtime",
	Help:      "How long has an action runner been up."}, []string{"label", "instance_id"})
