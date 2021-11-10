package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func (c *CommanderMetrics) initializeRollupLoopMetrics() {
	c.initializeBatchMetrics()
}

func (c *CommanderMetrics) initializeBatchMetrics() {
	buildAndSubmissionTimes := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: batchSubsystem,
			Name:      "build_and_submission_time_seconds",
			Help:      "Batch build and submission times",
			Buckets: []float64{
				0.1,
				0.2,
				0.3,
				0.4,
				0.5,
				1.0,
				2.0,
				3.0,
				4.0,
				5.0,
				10.0,
				15.0,
				20.0,
				25.0,
				50.0,
			},
		},
		[]string{"type"},
	)

	c.registry.MustRegister(
		buildAndSubmissionTimes,
	)

	c.BatchBuildAndSubmissionTimes = buildAndSubmissionTimes
}
