package components

import metrics "github.com/rcrowley/go-metrics"
import "time"

const (
	sampleSize      int     = 1028
	sampleAlpha     float64 = 0.015
	responseTimeKey string  = "ResponseTimeMetric"
	errorCodeKey    string  = "ErrorCodesGauge"
)

/**
TODO
Needed metrics are
- Gauge for # of bad requests
- Histogram and ExpSample for response time
- ?? for others
*/

// ResponseAggregator is a Singleton that captures and exposes the metric data
type ResponseAggregator interface {
	ListenAndAggregate(<-chan RequestResult)
}

type GoMetricBasedAggregator struct {
	config                *Config
	responseTimeHistogram metrics.Histogram
	errorCodesGauge       metrics.Gauge
}

// ListenAndAggregate sets up metrics and listens on the channel for completed requests
// aggregating accordingly
func (aggregator *GoMetricBasedAggregator) ListenAndAggregate(results <-chan RequestResult) {
	aggregator.setupMetrics()
	go aggregator.aggregateResults(results)
}

func (aggregator *GoMetricBasedAggregator) setupMetrics() {
	config := aggregator.config
	aggregator.responseTimeHistogram = getResponseTimMetricsHistogram(config)
	aggregator.errorCodesGauge = getErrorCodeGauge(config)
	// TODO figure out what life will look like for other metrics
}

func getResponseTimMetricsHistogram(config *Config) metrics.Histogram {
	s := metrics.NewExpDecaySample(sampleSize, sampleAlpha) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	metrics.Register(responseTimeKey, h)
	return h
}

func getErrorCodeGauge(config *Config) metrics.Gauge {
	g := metrics.NewGauge()
	metrics.Register(errorCodeKey, g)
	return g
}

func (aggregator *GoMetricBasedAggregator) aggregateResults(results <-chan RequestResult) {
	for result := range results {
		timeTaken := result.EndTime.Sub(result.StartTime) * time.Millisecond
		aggregator.responseTimeHistogram.Update(int64(timeTaken))
		aggregator.errorCodesGauge.Update(result.)
	}
}
