package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func (c *CommanderMetrics) initializeRollupLoopMetrics() {
	batchBuildAndSubmissionDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: rollupSubsystem,
			Name:      "batch_build_and_submission_duration_milliseconds",
			Help:      "Batch build and submission durations",
			Buckets: []float64{
				0.0,
				50.0,
				100.0,
				250.0,
				500.0,
				1000.0,
				2000.0,
				3000.0,
				4000.0,
				5000.0,
				6000.0,
				7000.0,
				8000.0,
				9000.0,
				10000.0,
				15000.0,
				20000.0,
			},
		},
		[]string{"type"},
	)

	commitmentBuildDurations := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: rollupSubsystem,
			Name:      "commitment_build_duration_milliseconds",
			Help:      "Commitment build durations",
			Buckets: []float64{
				0.0,
				25.0,
				50.0,
				75.0,
				100.0,
				150.0,
				200.0,
				250.0,
				300.0,
				350.0,
				400.0,
				450.0,
				500.0,
				600.0,
				700.0,
				800.0,
				900.0,
				1000.0,
			},
		},
		[]string{"type"},
	)

	c.registry.MustRegister(
		batchBuildAndSubmissionDuration,
		commitmentBuildDurations,
	)

	c.BatchBuildAndSubmissionTimes = batchBuildAndSubmissionDuration
	c.CommitmentBuildDuration = commitmentBuildDurations
}
